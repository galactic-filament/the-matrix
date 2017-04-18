package client

import (
	"strconv"

	"fmt"

	"errors"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/endpoint"
	"github.com/ihsw/the-matrix/app/repo"
	"github.com/ihsw/the-matrix/app/simpledocker"
)

func getContainerName(clientEndpoint endpoint.Endpoint, name string) string {
	return fmt.Sprintf("%s-%s-client", clientEndpoint.Name, name)
}

// DefaultAppPort - the expected port that the endpoint listens on
const DefaultAppPort = 80

// NewClient - creates a client to run against endpoints
func NewClient(clientRepo repo.Repo, clientNetwork *docker.Network) Client {
	return Client{
		Repo:    clientRepo,
		Network: clientNetwork,
	}
}

// Client - a repo for running against endpoints
type Client struct {
	repo.Repo
	Network   *docker.Network
	Container *docker.Container
}

// Run - runs this client against an endpoint
func (c Client) Run(clientEndpoint endpoint.Endpoint) (string, error) {
	// gathering the endpoint's ip address
	endpointHostIP, err := c.Client.GetContainerIP(c.Network, clientEndpoint.Container)
	if err != nil {
		return "", err
	}

	// validating that there is no client with this name
	containerName := getContainerName(clientEndpoint, c.Name)
	container, err := c.Client.GetContainer(containerName)
	if err != nil {
		return "", err
	}
	if container != nil {
		return "", errors.New("Client container already exists")
	}

	// creating the client container
	clientEnvVars := map[string]string{
		"APP_HOST": endpointHostIP.String(),
		"APP_PORT": strconv.Itoa(DefaultAppPort),
	}
	container, err = c.Client.CreateContainer(simpledocker.CreateContainerOptions{
		Name:    containerName,
		Image:   repo.GetImageName(c.Name),
		Network: c.Network,
		EnvVars: clientEnvVars,
	})
	if err != nil {
		return "", err
	}
	c.Container = container

	// running it out
	failed, err := c.Client.RunContainer(container, []string{})
	if err != nil {
		return "", err
	}
	if failed {
		containerOutput, err := c.Client.GetContainerLogs(c.Container)
		if err != nil {
			return "", err
		}

		return containerOutput, nil
	}

	return "", nil
}

// Clean - cleans a client's container
func (c Client) Clean() error {
	if err := c.Client.RemoveContainer(c.Container); err != nil {
		return err
	}

	return nil
}
