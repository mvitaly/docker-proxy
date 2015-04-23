package docker

import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
)

// Options for creating a new docker client
type Options struct {
	Address string
	CertPath string
}

// Docker client
type Docker struct {
	Client *docker.Client
}

// Container information
type Container struct {
	ID    string
	Image string
	Names []string
}

// Lister is an interface of to get APIContainers
type Lister interface {
	ListContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error)
}

// New creates a new instance of a docker client
func New(o *Options) (*Docker, error) {
	var client *docker.Client
	var err error
	if (o.CertPath == "") {
		client, err = docker.NewClient(o.Address)
	} else {
		ca := fmt.Sprintf("%s/ca.pem", o.CertPath)
		cert := fmt.Sprintf("%s/cert.pem", o.CertPath)
		key := fmt.Sprintf("%s/key.pem", o.CertPath)
		client, err = docker.NewTLSClient(o.Address, cert, key, ca)
	}
	if err != nil {
		return nil, err
	}
	return &Docker{
		Client: client,
	}, nil
}

// ContainersFilter are the options to filter based on containers
type ContainersFilter struct {
	ID string
}

// Containers returns a list of containers based on the filter options
func (d *Docker) Containers(f *ContainersFilter) ([]Container, error) {
	result := []Container{}
	if f == nil {
		f = &ContainersFilter{}
	}
	containers, err := d.Client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return nil, err
	}
	for _, container := range containers {
		if f.ID == "" {
			result = append(result, Container{
				ID:    container.ID,
				Image: container.Image,
				Names: container.Names,
			})
		} else {
			if container.ID == f.ID {
				result = append(result, Container{
					ID:    container.ID,
					Image: container.Image,
					Names: container.Names,
				})
			}
		}
	}
	return result, nil
}
