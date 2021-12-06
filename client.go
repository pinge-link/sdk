package pinge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/pinge-link/sdk/spec"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	ErrorAllGatesBusy = errors.New("all gates busy")
)

type connectOptions struct {
	Kind    int    `json:"kind"`
	Token   string `json:"token"`
	Service string `json:"service"`
	Private bool   `json:"private,omitempty"`
}

type Client struct {
	gateHost        string
	initHost        string
	accepter        chan net.Conn
	token           string
	private         bool
	serviceName     string
	customDomain    string
	topologyAddress string
	region          *TopologyRegion
	ctx             context.Context
	uri             string
}

type ClientOption func(*Client)

func WithTopologyAddress(address string) ClientOption {
	return func(c *Client) {
		c.topologyAddress = address
	}
}

func WithGateHost(host string) ClientOption {
	return func(c *Client) {
		c.gateHost = host
	}
}

func WithInitHost(host string) ClientOption {
	return func(c *Client) {
		c.initHost = host
	}
}

func WithPrivate() ClientOption {
	return func(c *Client) {
		c.private = true
	}
}

func WithCustomDomain(domain string) ClientOption {
	return func(c *Client) {
		c.customDomain = domain
	}
}

func InitClient(ctx context.Context, serviceName string, token string, options ...ClientOption) (*Client, error) {
	topologyDefault := "http://topology.pinge.dev:5004"

	c := Client{
		accepter:        make(chan net.Conn),
		token:           token,
		serviceName:     serviceName,
		topologyAddress: topologyDefault,
		ctx:             ctx,
	}

	for _, option := range options {
		option(&c)
	}

	topology, err := c.getTopology()
	if err != nil {
		return nil, fmt.Errorf("cannot get topology: %w", err)
	}

	region, err := c.selectRegion(topology)
	if err != nil {
		return nil, err
	}

	c.region = region

	c.gateHost = region.Gates[0].SecondaryAddress
	c.initHost = region.Gates[0].PrimaryAddress

	c.initPrimary()

	return &c, nil
}

func (c *Client) getTopology() (*TopologyConfig, error) {
	res, err := http.Get(c.topologyAddress)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var cfg TopologyConfig

	if err := json.NewDecoder(res.Body).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Client) selectRegion(topology *TopologyConfig) (*TopologyRegion, error) {
	var bestRegion *TopologyRegion
	var bestPingTime time.Duration

	for i, region := range topology.Regions {
		opts := []grpc.DialOption{
			grpc.WithInsecure(),
		}

		conn, err := grpc.Dial(region.PingHost, opts...)
		if err != nil {
			fmt.Printf("host %s for region %s, not available: %s\r\n", region.Id, region.PingHost, err)
			continue
		}

		pingerClient := spec.NewServiceClient(conn)

		startTime := time.Now()

		if _, err := pingerClient.Ping(context.Background(), &spec.PingRequestResponse{}); err != nil {
			fmt.Printf("host %s for region %s, not available: %s\r\n", region.Id, region.PingHost, err)
			continue
		}

		execTime := time.Since(startTime)

		if bestRegion == nil {
			bestRegion = &topology.Regions[i]
			bestPingTime = execTime
		} else if execTime < bestPingTime {
			bestRegion = &topology.Regions[i]
			bestPingTime = execTime
		}
	}

	if bestRegion == nil {
		return nil, fmt.Errorf("cannot find available gates")
	}

	return bestRegion, nil
}

func (c *Client) initPrimary() error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(c.initHost, opts...)
	if err != nil {
		return err
	}

	gateClient := spec.NewServiceClient(conn)

	stream, err := gateClient.Connect(c.ctx, &spec.ConnectRequest{
		Token:        c.token,
		ServiceName:  c.serviceName,
		Private:      c.private,
		CustomDomain: c.customDomain,
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				respStatus, ok := status.FromError(err)
				if ok {
					log.Fatal(respStatus.Message())
				}

				/*if strings.Contains(err.Error(), "cannot get token") {
					log.Fatalf("token %s is invalid \r\n", c.token)
				} else if strings.Contains(err.Error(), "must contain English letters and digits only") {
					log.Fatal(err)
				} else {
					fmt.Println(err)
				}*/

				select {
				case <-c.ctx.Done():
					return
				default:
				}

				if strings.Contains(err.Error(), "service exist") { // use grpc error
					if err := c.busyGate(); err != nil {
						log.Fatal(err)
						return
					}
				}
			reconnectLoop:
				for {
					select {
					case <-time.After(time.Second):
						fmt.Println("reconnect")
						if err := c.initPrimary(); err == nil {
							fmt.Println("reconnect error", err)
							break reconnectLoop
						}
					}
				}
				return
			}

			switch resp.Kind {
			case spec.Type_OPEN:
				secondConn, err := c.getConnection(c.serviceName, c.token)
				if err != nil {
					fmt.Println("cannot get connection")
					return
				}

				c.accepter <- secondConn
			case spec.Type_SET_INFO:
				c.uri = resp.ProjectUri

				fmt.Printf("Service URL: https://%s\r\n", resp.ProjectUri)
			}
		}
	}()

	return nil
}

func (c *Client) busyGate() error {
	var newGate *TopologyGate
	for i, gate := range c.region.Gates {
		if gate.PrimaryAddress == c.initHost {
			c.region.Gates[i].Busy = true
		} else if !gate.Busy {
			newGate = &c.region.Gates[i]
		}
	}

	if newGate == nil {
		return ErrorAllGatesBusy
	}

	c.gateHost = newGate.SecondaryAddress
	c.initHost = newGate.PrimaryAddress

	return nil
}

func (c *Client) Accept() (net.Conn, error) {
	select {
	case <-c.ctx.Done():
		return nil, fmt.Errorf("context deadline")
	default:
	}

	return <-c.accepter, nil
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) Addr() net.Addr {
	return nil
}

func (c *Client) getConnection(serviceName string, token string) (net.Conn, error) {
	conn, err := net.Dial("tcp", c.gateHost)
	if err != nil {
		return nil, err
	}

	kind := 3

	b, err := json.Marshal(connectOptions{
		Kind:    kind,
		Token:   token,
		Service: serviceName,
		Private: c.private,
	})
	if err != nil {
		return nil, err
	}

	b = append(b, 10)

	if _, err := conn.Write(b); err != nil {
		return nil, err
	}

	return conn, nil
}
