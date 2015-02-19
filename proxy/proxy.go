package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	dockerclient "github.com/fsouza/go-dockerclient"
	"github.com/xbudex/docker-proxy/docker"
)

// Options for making a new proxy
type Options struct {
	Docker docker.Lister
}

// New creates a new ReverseProxy instance
func New(o *Options) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			containers, _ := o.Docker.ListContainers(dockerclient.ListContainersOptions{})
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
							fmt.Printf("%s://%s\n", req.URL.Scheme, req.URL.Host)
							return
						}
					}
				}
				fmt.Printf("found nothing\n")
			}
		},

		ErrorLog: log.New(os.Stderr, "", log.Lshortfile),
	}
}
