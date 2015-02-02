package main

import (
	docker "github.com/fsouza/go-dockerclient"
)

type DockerOptions struct {
	Address string
}

type Docker struct {
	client *docker.Client
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
