package SimpleDocker

import (
	"bytes"
	"fmt"
	"github.com/fsouza/go-dockerclient"
)

// SimpleDocker - simplified interface to the docker client
type SimpleDocker struct {
	client *docker.Client
}

// NewSimpleDocker - creates a new SimpleDocker
func NewSimpleDocker(endpoint string) (SimpleDocker, error) {
	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		return SimpleDocker{}, err
	}

	return SimpleDocker{client: client}, nil
}

// CreateContainer - creates a container but doesn't start it up
func (s SimpleDocker) CreateContainer(name string, image string) (*docker.Container, error) {
	container, err := s.client.CreateContainer(docker.CreateContainerOptions{
		Name:   name,
		Config: &docker.Config{Image: image},
	})
	if err != nil {
		return nil, err
	}

	return container, nil
}

// GetContainerLogs - fetches the logs for a given container
func (s SimpleDocker) GetContainerLogs(container *docker.Container) (string, error) {
	var output bytes.Buffer
	err := s.client.Logs(docker.LogsOptions{
		Container:    container.ID,
		Stdout:       true,
		OutputStream: &output,
	})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// RunContainer - starts a container against an ApiServer container
func (s SimpleDocker) RunContainer(container *docker.Container, apiContainerID string) (bool, error) {
	err := s.client.StartContainer(container.ID, &docker.HostConfig{
		Links: []string{fmt.Sprintf("%s:ApiServer", apiContainerID)},
	})
	if err != nil {
		return false, err
	}

	status, err := s.client.WaitContainer(container.ID)
	if err != nil {
		return false, err
	}
	if status != 0 {
		return true, nil
	}

	return false, nil
}

// RemoveContainer - removes a container
func (s SimpleDocker) RemoveContainer(container *docker.Container) error {
	err := s.client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            container.ID,
		RemoveVolumes: true,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetContainer - fetches a container from id
func (s SimpleDocker) GetContainer(id string) (*docker.Container, error) {
	return s.client.InspectContainer(id)
}
