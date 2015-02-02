package main

import (
	"fmt"
	"net/http"

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

	proxy := getProxy(&proxyOptions{docker: client})

	http.ListenAndServe("localhost:8888", proxy)
}
