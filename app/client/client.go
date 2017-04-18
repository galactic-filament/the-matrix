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
func (c Client) Run(clientEndpoint endpoint.Endpoint) error {
	// gathering the endpoint's ip address
	endpointHostIP, err := c.Client.GetContainerIP(c.Network, clientEndpoint.Container)
	if err != nil {
		return err
	}
	if endpointHostIP == nil {
		return errors.New("Endpoint IP was nil")
	}

	// creating the client container
	clientEnvVars := map[string]string{
		"APP_HOST": endpointHostIP.String(),
		"APP_PORT": strconv.Itoa(DefaultAppPort),
	}
	container, err := c.Client.CreateContainer(simpledocker.CreateContainerOptions{
		Name:    getContainerName(clientEndpoint, c.Name),
		Image:   repo.GetImageName(c.Name),
		Network: c.Network,
		EnvVars: clientEnvVars,
	})
	c.Container = container

	// running it out
	failed, err := c.Client.RunContainer(container, []string{})
	if err != nil {
		return err
	}
	if failed {
		return errors.New("Client container failed")
	}

	return nil
}
