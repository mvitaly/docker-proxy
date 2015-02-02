package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	containers, _ := client.ListContainers(docker.ListContainersOptions{})
	for _, img := range containers {
		fmt.Println("ID: ", img.ID)
		fmt.Println("Image: ", img.Image)
		fmt.Println("Status:", img.Status)
		fmt.Println("Ports:", img.Ports)
		for _, port := range img.Ports {
			fmt.Println("\tPrivate:", port.PrivatePort)
			fmt.Println("\tPublic:", port.PublicPort)
			fmt.Println("\tType:", port.Type)
			fmt.Println("\tIP:", port.IP)
			fmt.Println("")
		}
		fmt.Println("Names:", img.Names)
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			containers, _ := client.ListContainers(docker.ListContainersOptions{})
			for _, container := range containers {
				imageParts := strings.Split(container.Image, ":")
				name := ""
				if len(imageParts) >= 1 {
					name = imageParts[0]
				}

				urlParts := strings.Split(req.Host, ":")
				host := urlParts[0]
				for _, subDomain := range strings.Split(host, ".") {
					if subDomain == name {
						for _, port := range container.Ports {
							if port.IP == "" || port.PublicPort == 0 {
								continue
							}
							req.URL.Scheme = "http"
							req.URL.Host = fmt.Sprintf("%s:%d", port.IP, port.PublicPort)
							fmt.Println(req.URL)
							return
						}
					}
				}
			}
		},

		FlushInterval: 0,

		ErrorLog: log.New(os.Stderr, "", log.Lshortfile),
	}

	http.ListenAndServe("localhost:8888", proxy)
}
