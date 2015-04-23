package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"github.com/xbudex/docker-proxy/docker"
	"github.com/xbudex/docker-proxy/proxy"
)

func mainAction(c *cli.Context) {
    fmt.Printf("Connecting to docker\n")
	client, err := docker.New(&docker.Options{Address: c.String("docker-host"), CertPath: c.String("docker-cert-path")})
	if err != nil {
		panic(err)
	}

	var targetPort = c.Int("target-port")
	if targetPort == 0 {
	    targetPort = c.Int("port")
	}

	proxy := proxy.New(&proxy.Options{Docker: client.Client, Port: targetPort, DefaultContainer: c.String("default-container-name")})
	bind := fmt.Sprintf("%s:%d", c.String("address"), c.Int("port"))

    fmt.Printf("Start listening to incoming connections on: %s\n", bind)
	panic(http.ListenAndServe(bind, proxy))
}

func main() {
	app := cli.NewApp()
	app.Name = "docker-proxy"
	app.Usage = "http proxy to a docker instance"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "docker-host, d",
			Value:  "unix:///var/run/docker.sock",
			Usage:  "address for docker socket",
			EnvVar: "DOCKER_HOST",
		},
		cli.StringFlag{
			Name:   "docker-cert-path, c",
			Value:  "",
			Usage:  "path to TLS certificates",
			EnvVar: "DOCKER_CERT_PATH",
		},
		cli.StringFlag{
		    Name:   "default-container-name, n",
			Usage:  "path to TLS certificates",
		},
		cli.IntFlag{
			Name:   "target-port, t",
			Value:  0,
			Usage:  "port to forward requests to, defaults to listen port",
		},
		cli.StringFlag{
			Name:   "address, a",
			Value:  "0.0.0.0",
			Usage:  "address to bind to",
			EnvVar: "DOCKER_PROXY_ADDRESS",
		},
		cli.IntFlag{
			Name:   "port, p",
			Value:  8080,
			Usage:  "port to bind to",
			EnvVar: "DOCKER_PROXY_PORT",
		},
	}
	app.Action = mainAction
	app.Run(os.Args)
}
