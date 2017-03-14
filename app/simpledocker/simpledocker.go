package simpledocker

import (
	"bytes"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

// SimpleDocker - simplified interface to the docker client
type SimpleDocker struct {
	client *docker.Client
}

// NewSimpleDocker - creates a new SimpleDocker
func NewSimpleDocker(dockerClient *docker.Client) SimpleDocker {
	return SimpleDocker{client: dockerClient}
}

// CreateContainer - creates a container but doesn't start it up
func (s SimpleDocker) CreateContainer(name string, image string, envVars []string) (*docker.Container, error) {
	container, err := s.client.CreateContainer(docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Image: image,
			Env:   envVars,
		},
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

// StartContainer - starts a container up
func (s SimpleDocker) StartContainer(container *docker.Container, links []string) error {
	log.WithFields(log.Fields{
		"container": container.Name,
		"links":     links,
	}).Info("Starting container with links")

	err := s.client.StartContainer(container.ID, &docker.HostConfig{Links: links})
	if err != nil {
		return err
	}

	return nil
}

// RunContainer - starts a container with links and waits for it to exit
func (s SimpleDocker) RunContainer(container *docker.Container, links []string) (bool, error) {
	err := s.StartContainer(container, links)
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

// StopContainer - stops a container
func (s SimpleDocker) StopContainer(container *docker.Container) error {
	return s.client.StopContainer(container.ID, 10)
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

// GetImage - fetches an image from id
func (s SimpleDocker) GetImage(id string) (*docker.Image, error) {
	return s.client.InspectImage(id)
}
