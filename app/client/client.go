package client

import (
	"errors"
	"strconv"

	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/endpoint"
	"github.com/ihsw/the-matrix/app/repo"
	"github.com/ihsw/the-matrix/app/simpledocker"
)

func getContainerName(clientEndpoint endpoint.Endpoint, name string) string {
	return fmt.Sprintf("%s-%s-client", clientEndpoint.Name, name)
}

// ErrClientFailed - returned when running a client fails
var ErrClientFailed = errors.New("Client failed")

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
	Network *docker.Network
}

// Run - runs this client against an endpoint
func (c Client) Run(clientEndpoint endpoint.Endpoint) (*docker.Container, error) {
	// gathering the endpoint's ip address
	endpointHostIP, err := c.Client.GetContainerIP(c.Network, clientEndpoint.Container)
	if err != nil {
		return nil, err
	}

	// creating the client container
	clientEnvVars := map[string]string{
		"API_HOST": endpointHostIP.String(),
		"API_PORT": strconv.Itoa(DefaultAppPort),
	}
	containerName := getContainerName(clientEndpoint, c.Name)
	createContainerOpts := simpledocker.CreateContainerOptions{
		Name:    containerName,
		Image:   repo.GetImageName(c.Name),
		Network: c.Network,
		EnvVars: clientEnvVars,
	}
	container, err := c.Client.CreateContainer(createContainerOpts)
	if err != nil {
		// failing on real error
		if err != docker.ErrContainerAlreadyExists {
			return nil, err
		}

		// removing the existing client container
		container, err := c.Client.GetContainer(containerName)
		if err != nil {
			return nil, err
		}
		if err := c.Client.RemoveContainer(container); err != nil {
			return nil, err
		}

		// creating a new client container
		container, err = c.Client.CreateContainer(createContainerOpts)
		if err != nil {
			return nil, err
		}
	}

	// running it out
	failed, err := c.Client.RunContainer(container, []string{})
	if err != nil {
		return nil, err
	}
	if failed {
		return container, ErrClientFailed
	}

	return container, nil
}

// Clean - cleans a client's container
func (c Client) Clean(container *docker.Container) error {
	if err := c.Client.RemoveContainer(container); err != nil {
		return err
	}

	return nil
}
