package pinge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func DockerInit(token string, initHost string) error {
	return getContainers("/var/run/docker.sock", token, initHost)
}

func watchContainer(ctx context.Context, dockerSockPath string, containerName string) (chan *DockerContainerEvent, error) {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", dockerSockPath)
			},
		},
	}

	var filter string

	if containerName != "" {
		filter = `filters={"container":["` + containerName + `"]}`
	}

	res, err := httpc.Get("http://unix/v1.24/events?" + filter)
	if err != nil {
		return nil, err
	}

	ch := make(chan *DockerContainerEvent)

	go func() {
		defer res.Body.Close()

		dec := json.NewDecoder(res.Body)

		for {
			var event DockerContainerEvent

			select {
			case <-ctx.Done():
				fmt.Println("stop listen container events")
				return
			default:
			}

			if err := dec.Decode(&event); err != nil {
				if err == io.EOF {
					// all done
					break
				}

				if err != nil {
					return
				}
			}

			ch <- &event
		}
	}()

	return ch, nil
}

func startContainer(dockerSockPath string, token string, container DockerContainer, initHost string) error {
	pingeService := container.Labels["pingeService"]
	pingePort := container.Labels["pingePort"]
	pingeContainerPort := container.Labels["pingeContainerPort"]

	if pingePort == "" && pingeContainerPort == "" {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	var host string
	var port string

	if pingePort != "" {
		host = "localhost"
		port = pingePort
	} else {
		host = container.NetworkSettings.Networks.Bridge.IPAddress
		port = pingeContainerPort
	}

	updateEvents, err := watchContainer(ctx, dockerSockPath, container.ID)
	if err != nil {
		cancel()
		return err
	}

	g := errgroup.Group{}

	g.Go(func() error {
		for event := range updateEvents {
			if event.Action == "stop" {
				fmt.Println("stop container")
				cancel()
				return nil
			}
		}

		return nil
	})

	g.Go(func() error {
		fmt.Println("start service", pingeService, host, port)
		var options []ClientOption

		if initHost != "" {
			options = append(options, WithTopologyAddress(initHost))
		}

		return InitService(ctx, pingeService, token, host, port, options)
	})

	return g.Wait()
}

func getContainers(dockerSockPath string, token string, initHost string) error {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", dockerSockPath)
			},
		},
	}

	filter := `filters={"label":["pingeService"]}`
	res, err := httpc.Get("http://unix/v1.24/containers/json?all=1&before=8dfafdbc3a40&size=1&" + filter)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("mismatch statuses: %v", res.StatusCode)
	}

	var containers []DockerContainer

	if err := json.NewDecoder(res.Body).Decode(&containers); err != nil {
		return err
	}

	g := errgroup.Group{}

	for _, container := range containers {
		if container.GetState() == "running" {
			go startContainer(dockerSockPath, token, container, initHost)
		}
	}

	g.Go(func() error {
		ch, err := watchContainer(context.Background(), dockerSockPath, "")
		if err != nil {
			return err
		}

		fmt.Println("listen containers")

		for event := range ch {
			fmt.Println("[[[[[[")
			if event.Action == "start" {
				fmt.Println("container listener receive event", event.Action, event.Actor.ID)

				res, err := httpc.Get("http://unix/v1.24/containers/" + event.Actor.ID + "/json")
				if err != nil {
					fmt.Println(err)
					return err
				}

				var container DockerContainer

				if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
					res.Body.Close()
					fmt.Println(err)
					return err
				}

				container.Labels = container.Config.Labels

				go startContainer(dockerSockPath, token, container, initHost)
			}
		}

		return nil
	})

	return g.Wait()
}

type DockerContainer struct {
	ID      string      `json:"Id"`
	Names   []string    `json:"Names"`
	Image   string      `json:"Image"`
	ImageID string      `json:"ImageID"`
	State   interface{} `json:"state"`
	Command string      `json:"Command"`
	Ports   []struct {
		PrivatePort int    `json:"PrivatePort"`
		Type        string `json:"Type"`
	} `json:"Ports"`
	SizeRw     int               `json:"SizeRw"`
	SizeRootFs int               `json:"SizeRootFs"`
	Labels     map[string]string `json:"Labels"`
	Status     string            `json:"Status"`
	HostConfig struct {
		NetworkMode string `json:"NetworkMode"`
	} `json:"HostConfig"`
	Config struct {
		Labels map[string]string `json:"Labels"`
	}
	NetworkSettings struct {
		Networks struct {
			Bridge struct {
				IPAMConfig          interface{} `json:"IPAMConfig"`
				Links               interface{} `json:"Links"`
				Aliases             interface{} `json:"Aliases"`
				NetworkID           string      `json:"NetworkID"`
				EndpointID          string      `json:"EndpointID"`
				Gateway             string      `json:"Gateway"`
				IPAddress           string      `json:"IPAddress"`
				IPPrefixLen         int         `json:"IPPrefixLen"`
				IPv6Gateway         string      `json:"IPv6Gateway"`
				GlobalIPv6Address   string      `json:"GlobalIPv6Address"`
				GlobalIPv6PrefixLen int         `json:"GlobalIPv6PrefixLen"`
				MacAddress          string      `json:"MacAddress"`
				DriverOpts          interface{} `json:"DriverOpts"`
			} `json:"bridge"`
		} `json:"Networks"`
	} `json:"NetworkSettings"`
	Mounts []interface{} `json:"Mounts"`
}

func (d *DockerContainer) GetState() string {
	switch x := d.State.(type) {
	case string:
		return x
	}

	return ""
}

type DockerContainerEvent struct {
	Status string `json:"status"`
	ID     string `json:"id"`
	From   string `json:"from"`
	Type   string `json:"Type"`
	Action string `json:"Action"`
	Actor  struct {
		ID         string            `json:"ID"`
		Attributes map[string]string `json:"Attributes"`
	} `json:"Actor"`
	Time     int   `json:"time"`
	TimeNano int64 `json:"timeNano"`
}

type DockerContainerItem struct {
	ID      string   `json:"Id"`
	Names   []string `json:"Names"`
	Image   string   `json:"Image"`
	ImageID string   `json:"ImageID"`
	Command string   `json:"Command"`
	Ports   []struct {
		PrivatePort int    `json:"PrivatePort"`
		Type        string `json:"Type"`
	} `json:"Ports"`
	SizeRw     int               `json:"SizeRw"`
	SizeRootFs int               `json:"SizeRootFs"`
	Labels     map[string]string `json:"Labels"`
	Status     string            `json:"Status"`
	HostConfig struct {
		NetworkMode string `json:"NetworkMode"`
	} `json:"HostConfig"`
	Config struct {
		Labels map[string]string `json:"Labels"`
	}
	NetworkSettings struct {
		Networks struct {
			Bridge struct {
				IPAMConfig          interface{} `json:"IPAMConfig"`
				Links               interface{} `json:"Links"`
				Aliases             interface{} `json:"Aliases"`
				NetworkID           string      `json:"NetworkID"`
				EndpointID          string      `json:"EndpointID"`
				Gateway             string      `json:"Gateway"`
				IPAddress           string      `json:"IPAddress"`
				IPPrefixLen         int         `json:"IPPrefixLen"`
				IPv6Gateway         string      `json:"IPv6Gateway"`
				GlobalIPv6Address   string      `json:"GlobalIPv6Address"`
				GlobalIPv6PrefixLen int         `json:"GlobalIPv6PrefixLen"`
				MacAddress          string      `json:"MacAddress"`
				DriverOpts          interface{} `json:"DriverOpts"`
			} `json:"bridge"`
		} `json:"Networks"`
	} `json:"NetworkSettings"`
	Mounts []interface{} `json:"Mounts"`
}
