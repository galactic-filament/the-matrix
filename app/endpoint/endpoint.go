package endpoint

import (
	"errors"
	"fmt"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/repo"
	"github.com/ihsw/the-matrix/app/resource"
)

// NewEndpoint - creates a new endpoint for a client to consume
func NewEndpoint(repo repo.Repo, resources resource.Resources) (Endpoint, error) {
	endpoint := Endpoint{repo, nil}

	// creating an endpoint container
	container, err := repo.Client.CreateContainer(
		fmt.Sprintf("%s-endpoint", endpoint.Name),
		fmt.Sprintf("ihsw/%s", endpoint.Name),
		resources.GetEnvVarsList(),
	)
	if err != nil {
		return Endpoint{}, err
	}
	endpoint.Container = container

	// starting it up with links to the provided resources
	if err := endpoint.Client.StartContainer(container, resources.GetLinkLineList()); err != nil {
		return Endpoint{}, err
	}

	// waiting for the endpoint to start up
	time.Sleep(10 * time.Second)

	return endpoint, nil
}

// Endpoint - a container ran against an Endpoint
type Endpoint struct {
	repo.Repo
	Container *docker.Container
}

// Clean - stops and removes an Endpoint's container
func (e Endpoint) Clean() error {
	if e.Container == nil {
		return errors.New("Endpoint container was nil")
	}

	isRunning, err := e.Client.IsRunning(e.Container)
	if err != nil {
		return err
	}
	if isRunning {
		if err := e.Client.StopContainer(e.Container); err != nil {
			return err
		}
	}

	if err := e.Client.RemoveContainer(e.Container); err != nil {
		return err
	}

	return nil
}
