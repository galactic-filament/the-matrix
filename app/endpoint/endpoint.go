package endpoint

import (
	"errors"
	"fmt"

	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/galactic-filament/the-matrix/app/repo"
	"github.com/galactic-filament/the-matrix/app/resource"
	"github.com/galactic-filament/the-matrix/app/simpledocker"
	"github.com/galactic-filament/the-matrix/app/util"
)

func getContainerName(name string) (string, error) {
	name, err := util.GetPrefixedUUID(fmt.Sprintf("%s-endpoint", name))
	if err != nil {
		return "", err
	}

	return name, nil
}

// NewEndpoint - creates a new endpoint for a client to consume
func NewEndpoint(endpointRepo repo.Repo, network *docker.Network, resources resource.Resources) (Endpoint, error) {
	e := Endpoint{
		Repo:    endpointRepo,
		Network: network,
	}

	// gathering the host addresses of resources
	endpointEnvVars := map[string]string{}
	for _, resource := range resources.Values {
		ip, err := resource.GetContainerIP()
		if err != nil {
			return Endpoint{}, err
		}

		endpointEnvVars[fmt.Sprintf("%s_HOST", resource.EndpointLabel)] = ip.String()
	}

	// creating an endpoint container
	containerName, err := getContainerName(e.Name)
	if err != nil {
		return Endpoint{}, err
	}
	container, err := endpointRepo.Client.CreateContainer(simpledocker.CreateContainerOptions{
		Name:    containerName,
		Image:   repo.GetImageName(e.Name),
		Network: network,
		EnvVars: endpointEnvVars,
	})
	if err != nil {
		return Endpoint{}, err
	}
	e.Container = container

	// starting it up
	if err := e.Client.StartContainer(container, []string{}); err != nil {
		return Endpoint{}, err
	}

	// verifying that it is running
	isRunning, err := e.Client.IsRunning(e.Container)
	if err != nil {
		return Endpoint{}, err
	}
	if !isRunning {
		return Endpoint{}, errors.New("Endpoint was not up")
	}

	// sleeping to ensure this endpoint finishes booting up
	time.Sleep(util.PostDockerActionDelayInSeconds * time.Second)

	return e, nil
}

// Endpoint - a container ran against an Endpoint
type Endpoint struct {
	repo.Repo
	Network   *docker.Network
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
