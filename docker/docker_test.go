package docker

import (
	"fmt"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
)

type mockClient struct {
	containers []docker.APIContainers
	error      error
}

func (c *mockClient) ListContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error) {
	return c.containers, c.error
}

var container1 docker.APIContainers = docker.APIContainers{
	ID:     "01234567890a",
	Image:  "some-image:tag",
	Status: "Up Time",
	Ports: []docker.APIPort{
		docker.APIPort{
			PrivatePort: 12345,
			PublicPort:  78901,
			Type:        "TCP",
			IP:          "127.0.0.1",
		},
	},
	SizeRw:     1234567890123456789,
	SizeRootFs: 1234567890123456789,
	Names:      []string{"suspicious_yalow", "some-image"},
}

var container2 docker.APIContainers = docker.APIContainers{
	ID:     "bcdef0123456",
	Image:  "some-image:tag",
	Status: "Up Time",
	Ports: []docker.APIPort{
		docker.APIPort{
			PrivatePort: 12345,
			PublicPort:  78901,
			Type:        "TCP",
			IP:          "127.0.0.1",
		},
	},
	SizeRw:     1234567890123456789,
	SizeRootFs: 1234567890123456789,
	Names:      []string{"suspicious_yalow", "some-image"},
}

func containersMatch(t *testing.T, container Container, api docker.APIContainers) (bool, string) {
	if container.ID != api.ID {
		return false, fmt.Sprintf(
			"ID did not match. Got %s, expected %s.",
			container.ID,
			api.ID,
		)
	}
	if container.Image != api.Image {
		return false, fmt.Sprintf(
			"Image did not match. Got %s, expected %s.",
			container.Image,
			api.Image,
		)
	}
	if len(container.Names) != len(api.Names) {
		return false, fmt.Sprintf(
			"Has different number of names. Got %d, expected %d.",
			len(container.Names),
			len(api.Names),
		)
	}
	for index := range container.Names {
		if container.Names[index] != api.Names[index] {
			return false, fmt.Sprintf(
				"Has different number of names. Got %d, expected %d.",
				container.Names[index],
				api.Names[index],
			)
		}
	}
	return true, ""
}

func TestDockerGetContainer(t *testing.T) {
	sampleContainers := []docker.APIContainers{container1, container2}
	dockerClient := &Docker{Client: &mockClient{containers: sampleContainers}}

	containers, err := dockerClient.Containers(nil)
	if err != nil {
		t.Errorf("Got unexpected eror. %s", err)
	}
	count := 0
	for index, container := range containers {
		if ok, err := containersMatch(t, container, sampleContainers[index]); !ok {
			t.Errorf(err)
		}
		count += 1
	}
	if count != len(containers) {
		t.Errorf("Did not get the right number of containers. Got %d, expected %d.", count, len(containers))
	}
}

func TestDockerGetContainerByName(t *testing.T) {
	sampleContainers := []docker.APIContainers{container1, container2}
	dockerClient := &Docker{Client: &mockClient{containers: sampleContainers}}

	containers, err := dockerClient.Containers(&ContainersFilter{
		ID: "bcdef0123456",
	})
	if err != nil {
		t.Errorf("Got unexpected eror. %s", err)
	}

	if ok, err := containersMatch(t, containers[0], sampleContainers[1]); !ok {
		t.Errorf(err)
	}

	if 1 != len(containers) {
		t.Errorf("Did not get the right number of containers. Got %d, expected %d.", len(containers), 1)
	}
}
