package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	dockerclient "github.com/fsouza/go-dockerclient"
)

// Options for making a new proxy
type Options struct {
	Docker *dockerclient.Client
	Port int
	DefaultContainer string
}

// GetContainerHost is used to get the proxied host
func GetContainerHost(container *dockerclient.Container, host string, port string) string {
    var containerHost string
    for containerPort, containerHostPorts := range container.NetworkSettings.Ports {
        if string(containerPort) == port + "/tcp" {
            if containerHostPorts == nil {
                ip := container.NetworkSettings.IPAddress
                containerHost = fmt.Sprintf("%s:%s", ip, port)
            } else {
                containerHostPort := containerHostPorts[0]
                if containerHostPort.HostIP == "0.0.0.0" {
                    containerHost = fmt.Sprintf("%s:%s", host, containerHostPort.HostPort)
                } else {
                    containerHost = fmt.Sprintf("%s:%s", containerHostPort.HostIP, containerHostPort.HostPort)
                }
            }
        }
    }
    return containerHost
}

// New creates a new ReverseProxy instance
func New(o *Options) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
            urlParts := strings.Split(req.Host, ":")
            host := urlParts[0]
            port := fmt.Sprintf("%d", o.Port)

            fmt.Printf("Initial request to %s\n", host)

			containers, _ := o.Docker.ListContainers(dockerclient.ListContainersOptions{})
			LookForContainerByName:
                for _, container := range containers {

                    // Containers have multiple names, look for the one that is not a link name (only one slash in the name)
                    for _, containerName := range container.Names {
                        nameParts := strings.Split(containerName, "/")
                        if len(nameParts) > 2 {
                            continue
                        }

                        name := nameParts[1]

                        fmt.Printf("Check matching container: %s\n", name)

                        for _, subDomain := range strings.Split(host, ".") {
                            if subDomain == name {
                                fmt.Printf("Name match to container: %s\n", name)
                                inspectedContainer, _ := o.Docker.InspectContainer(container.ID)

                                containerHost := GetContainerHost(inspectedContainer, host, port)
                                if containerHost != "" {
                                    req.URL.Host = containerHost
                                    req.URL.Scheme = "http"
                                    fmt.Printf("Proxy request to %s://%s\n", req.URL.Scheme, req.URL.Host)

                                    return
                                }

                                fmt.Printf("No port match: %s\n", name)

                                break LookForContainerByName
                            }
                        }
                    }
                }

            fmt.Printf("Found no matching container\n")

            if o.DefaultContainer != "" {
                fmt.Printf("Using default container\n")
                inspectedContainer, _ := o.Docker.InspectContainer(o.DefaultContainer)

                containerHost := GetContainerHost(inspectedContainer, host, port)

                if containerHost != "" {
                    req.URL.Host = containerHost
                    req.URL.Scheme = "http"
                    fmt.Printf("Proxy request to %s://%s\n", req.URL.Scheme, req.URL.Host)

                    return
                }
            }

            fmt.Printf("Returning error\n")
            req.URL.Host = ""
		},

		ErrorLog: log.New(os.Stderr, "", log.Lshortfile),
	}
}
