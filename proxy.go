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

type proxyOptions struct {
	docker Lister
}

func getProxy(o *proxyOptions) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{

		Director: func(req *http.Request) {
			containers, _ := o.docker.ListContainers(docker.ListContainersOptions{})
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
							hostParts := strings.Split(req.Host, ":")
							req.URL.Host = fmt.Sprintf("%s:%d", hostParts[0], port.PublicPort)
							req.URL.Scheme = "http"
							return
						}
					}
				}
			}
		},

		ErrorLog: log.New(os.Stderr, "", log.Lshortfile),
	}
}
