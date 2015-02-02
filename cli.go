package main

import (
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
	http.ListenAndServe(c.String("address"), proxy)
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
			Value:  "localhost:8888",
			Usage:  "address to bind to",
			EnvVar: "DOCKER_PROXY_ADDRESS",
		},
	}
	app.Action = mainAction
	app.Run(os.Args)

}
