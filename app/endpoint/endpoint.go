package endpoint

import (
	"fmt"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/repo"
	"github.com/ihsw/the-matrix/app/resource"
)

func newEndpoint(repo repo.Repo, resources resource.Resources) (Endpoint, error) {
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

	// starting it up with links to the provided resources
	if err := endpoint.Client.StartContainer(container, resources.GetLinkLineList()); err != nil {
		return Endpoint{}, err
	}
	endpoint.Container = container

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
func (e Endpoint) Clean(prevErr error) error {
	if err := e.Client.StopContainer(e.Container); err != nil {
		return err
	}

	if err := e.Client.RemoveContainer(e.Container); err != nil {
		return err
	}

	return prevErr
}
