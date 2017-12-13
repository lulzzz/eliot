package main

import (
	"os"

	"github.com/ernoaapa/eliot/cmd"
	"github.com/ernoaapa/eliot/pkg/api"
	"github.com/ernoaapa/eliot/pkg/controller"
	"github.com/ernoaapa/eliot/pkg/device"
	"github.com/ernoaapa/eliot/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/thejerf/suture"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "eliotd"
	app.Usage = "Daemon which contains all Eliot, for example GRPC API for the CLI client"
	app.UsageText = `eliotd [arguments...]

	 # By default listen port 5000
	 eliotd
	
	 # Listen custom port
	 eliotd --listen 0.0.0.0:5001`
	app.Description = `API for create/update/delete the containers and a way to connect into the containers.`
	app.Flags = append([]cli.Flag{
		cli.StringFlag{
			Name:   "containerd",
			Usage:  "containerd socket path for containerd's GRPC server",
			EnvVar: "ELLIOT_CONTAINERD",
			Value:  "/run/containerd/containerd.sock",
		},
		cli.StringFlag{
			Name:   "listen",
			Usage:  "GRPC host:port what to listen for client connections",
			EnvVar: "ELLIOT_LISTEN",
			Value:  "localhost:5000",
		},
	}, cmd.GlobalFlags...)
	app.Version = version.VERSION
	app.Before = cmd.GlobalBefore

	app.Action = func(clicontext *cli.Context) error {
		resolver := device.NewResolver(cmd.GetLabels(clicontext))
		device := resolver.GetInfo()
		client := cmd.GetRuntimeClient(clicontext, device.Hostname)
		listen := clicontext.String("listen")

		supervisor := suture.NewSimple("eliotd")
		supervisor.Add(api.NewServer(listen, client))
		supervisor.Add(controller.NewLifecycle(client))

		supervisor.Serve()

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
