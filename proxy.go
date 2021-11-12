package pinge

import (
	"context"
	"fmt"
	"io"
	"net"

	"golang.org/x/sync/errgroup"
)

func InitService(ctx context.Context, serviceName string, token string, host string, port string, options []ClientOption) error {
	client, err := InitClient(ctx, serviceName, token, options...)
	if err != nil {
		return err
	}

	for {
		conn, err := client.Accept()
		if err != nil {
			return err
		}

		handler := func() error {
			defer conn.Close()

			localConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
			if err != nil {
				return err
			}

			defer localConn.Close()

			g := errgroup.Group{}

			g.Go(func() error {
				io.Copy(conn, localConn)
				return nil
			})

			g.Go(func() error {
				io.Copy(localConn, conn)
				return nil
			})

			return g.Wait()
		}

		go func() {
			if err := handler(); err != nil {
				fmt.Println("close connection with error", host, port, err)
			}
		}()
	}
}
