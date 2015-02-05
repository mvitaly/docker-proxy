package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
)

func mainAction(c *cli.Context) {
	docker, err := NewDocker(&DockerOptions{Address: c.String("docker-host")})
	if err != nil {
		panic(err)
	}

	proxy := getProxy(&proxyOptions{docker: docker.client})
	bind := fmt.Sprintf("%s:%d", c.String("address"), c.Int("port"))
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
			Name:   "address, a",
			Value:  "0.0.0.0",
			Usage:  "address to bind to",
			EnvVar: "DOCKER_PROXY_ADDRESS",
		},
		cli.IntFlag{
			Name:   "port, p",
			Value:  8080,
			Usage:  "address to bind to",
			EnvVar: "DOCKER_PROXY_PORT",
		},
	}
	app.Action = mainAction
	app.Run(os.Args)

}
