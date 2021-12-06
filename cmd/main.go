package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	client "github.com/pinge-link/sdk"
)

func main() {
	port := flag.String("port", "", "specify your application port")
	host := flag.String("gate", "", "specify gate host")
	initHost := flag.String("init-host", "", "specify init host")
	serviceName := flag.String("service-name", "", "specity service name")
	token := flag.String("token", "", "specity token for pinge.link")
	command := flag.String("command", "", "specify command for run")
	private := flag.Bool("private", false, "access to service by token")
	docker := flag.Bool("docker", false, "scan docker containers and pinge labels")

	flag.Parse()

	fmt.Println(*command)

	if *token == "" {
		*token = os.Getenv("PINGE_TOKEN")
		if *token == "" {
			log.Fatal("token is empty")
		}
	}

	if *docker == true {
		if err := client.DockerInit(*token, *initHost); err != nil {
			log.Fatal(err)
		}

		return
	}

	if *serviceName == "" {
		*serviceName = os.Getenv("PINGE_SERVICE_NAME")
		if *serviceName == "" {
			log.Fatal("service name is empty")
		}
	}

	if *port == "" {
		*port = os.Getenv("PINGE_PORT")
		if *port == "" {
			log.Fatal("port is empty")
		}
	}

	if *initHost == "" {
		*initHost = os.Getenv("PINGE_TOPOLOGY_HOST")
	}

	var options []client.ClientOption

	if *host != "" {
		options = append(options, client.WithGateHost(*host))
	}

	if *private {
		options = append(options, client.WithPrivate())
	}

	if *initHost != "" {
		options = append(options, client.WithTopologyAddress(*initHost))
	}

	if *command != "" {
		go execCommand(*command)
	}

	if err := client.InitService(context.Background(), *serviceName, *token, "localhost", *port, options); err != nil {
		log.Fatal(err)
	}
}

func execCommand(commandString string) {
	parts := strings.Split(commandString, " ")

	command := exec.Command(parts[0], parts[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
