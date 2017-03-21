package simpledocker

import (
	"bytes"

	"github.com/fsouza/go-dockerclient"
)

// Client - simplified interface to the docker client
type Client struct {
	dockerClient *docker.Client
}

// NewClient - creates a new SimpleDocker
func NewClient(dockerClient *docker.Client) Client {
	return Client{dockerClient: dockerClient}
}

// CreateContainer - creates a container but doesn't start it up
func (c Client) CreateContainer(name string, image string, envVars []string) (*docker.Container, error) {
	container, err := c.dockerClient.CreateContainer(docker.CreateContainerOptions{
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
func (c Client) GetContainerLogs(container *docker.Container) (string, error) {
	var output bytes.Buffer
	err := c.dockerClient.Logs(docker.LogsOptions{
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
func (c Client) StartContainer(container *docker.Container, links []string) error {
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

// GetImage - fetches an image from id
func (c Client) GetImage(id string) (*docker.Image, error) {
	return c.dockerClient.InspectImage(id)
}

// HasImage - validates whether an image exists
func (c Client) HasImage(id string) (bool, error) {
	_, err := c.GetImage(id)
	if err != nil {
		if err != docker.ErrNoSuchImage {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// PullImage - pulls an image by repository and tag
func (c Client) PullImage(repository string, tag string) error {
	return c.dockerClient.PullImage(
		docker.PullImageOptions{Repository: repository, Tag: tag},
		docker.AuthConfiguration{},
	)
}

// RemoveImage - removes an image by id
func (c Client) RemoveImage(imageID string) error { return c.dockerClient.RemoveImage(imageID) }
