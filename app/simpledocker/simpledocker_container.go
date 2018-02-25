package simpledocker

import (
	"bytes"
	"errors"
	"fmt"

	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/galactic-filament/the-matrix/app/util"
)

// CreateContainer - creates a container but doesn't start it up
func (c Client) CreateContainer(opts CreateContainerOptions) (*docker.Container, error) {
	container, err := c.dockerClient.CreateContainer(opts.toDockerOpts())
	if err != nil {
		return nil, err
	}

	time.Sleep(util.PostDockerActionDelayInSeconds * time.Second)

	return container, nil
}

// GetContainerLogs - fetches the logs for a given container
func (c Client) GetContainerLogs(container *docker.Container) (string, error) {
	var (
		stdOut bytes.Buffer
		stdErr bytes.Buffer
	)
	err := c.dockerClient.Logs(docker.LogsOptions{
		Container:    container.ID,
		Stdout:       true,
		OutputStream: &stdOut,
		Stderr:       true,
		ErrorStream:  &stdErr,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", stdOut.String(), stdErr.String()), nil
}

// StartContainer - starts a container up
func (c Client) StartContainer(container *docker.Container, links []string) error {
	if container == nil {
		return errors.New("Cannot start nil container")
	}

	err := c.dockerClient.StartContainer(container.ID, &docker.HostConfig{Links: links})
	if err != nil {
		return err
	}

	return nil
}

// RunContainer - starts a container with links and waits for it to exit
func (c Client) RunContainer(container *docker.Container, links []string) (bool, error) {
	err := c.StartContainer(container, links)
	if err != nil {
		return false, err
	}

	status, err := c.dockerClient.WaitContainer(container.ID)
	if err != nil {
		return false, err
	}
	if status != 0 {
		return true, nil
	}

	return false, nil
}

// StopContainer - stops a container
func (c Client) StopContainer(container *docker.Container) error {
	return c.dockerClient.StopContainer(container.ID, 10)
}

// RemoveContainer - removes a container
func (c Client) RemoveContainer(container *docker.Container) error {
	if container == nil {
		return errors.New("Cannot remove nil container")
	}

	err := c.dockerClient.RemoveContainer(docker.RemoveContainerOptions{
		ID:            container.ID,
		RemoveVolumes: true,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetContainer - fetches a container from id
func (c Client) GetContainer(id string) (*docker.Container, error) {
	return c.dockerClient.InspectContainer(id)
}

// IsRunning - checks whether a container is running
func (c Client) IsRunning(container *docker.Container) (bool, error) {
	container, err := c.dockerClient.InspectContainer(container.ID)
	if err != nil {
		return false, err
	}

	return container.State.Running, nil
}

// CreateContainerOptions - options for CreateContainer
type CreateContainerOptions struct {
	Name    string
	Image   string
	Network *docker.Network
	EnvVars map[string]string
}

func (c CreateContainerOptions) toDockerOpts() docker.CreateContainerOptions {
	envVars := []string{}
	for k, v := range c.EnvVars {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}

	var networkingConfig *docker.NetworkingConfig
	if c.Network != nil {
		networkingConfig = &docker.NetworkingConfig{
			EndpointsConfig: map[string]*docker.EndpointConfig{
				c.Network.Name: &docker.EndpointConfig{NetworkID: c.Network.ID},
			},
		}
	}

	return docker.CreateContainerOptions{
		Name:             c.Name,
		Config:           &docker.Config{Image: c.Image, Env: envVars},
		NetworkingConfig: networkingConfig,
	}
}

// GetContainersByImageID - fetches all containers on a particular imageID and tag
func (c Client) GetContainersByImageID(imageID string, tag string) ([]*docker.Container, error) {
	apiContainers, err := c.dockerClient.ListContainers(docker.ListContainersOptions{
		All:     true,
		Filters: map[string][]string{"ancestor": []string{fmt.Sprintf("%s:%s", imageID, tag)}},
	})
	if err != nil {
		return []*docker.Container{}, err
	}

	containers := []*docker.Container{}
	for _, apiContainer := range apiContainers {
		container, err := c.GetContainer(apiContainer.ID)
		if err != nil {
			return []*docker.Container{}, err
		}

		containers = append(containers, container)
	}

	return containers, nil
}
