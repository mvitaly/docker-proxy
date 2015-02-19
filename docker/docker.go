package docker

import (
	docker "github.com/fsouza/go-dockerclient"
)

// Options for creating a new docker client
type Options struct {
	Address string
}

// Docker client
type Docker struct {
	Client Lister
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
	client, err := docker.NewClient(o.Address)
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
