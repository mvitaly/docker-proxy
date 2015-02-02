package main

import (
	docker "github.com/fsouza/go-dockerclient"
)

type DockerOptions struct {
	Address string
}

type Docker struct {
	client Lister
}

type Container struct {
	ID    string
	Image string
	Names []string
}

type Lister interface {
	ListContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error)
}

func NewDocker(o *DockerOptions) (*Docker, error) {
	client, err := docker.NewClient(o.Address)
	if err != nil {
		return nil, err
	}
	return &Docker{
		client: client,
	}, nil
}

type ContainersFilter struct {
	ID string
}

func (d *Docker) Containers(f *ContainersFilter) ([]Container, error) {
	result := []Container{}
	if f == nil {
		f = &ContainersFilter{}
	}
	containers, err := d.client.ListContainers(docker.ListContainersOptions{})
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
