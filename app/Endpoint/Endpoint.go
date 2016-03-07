package Endpoint

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/Repo"
	"github.com/ihsw/the-matrix/app/Resource"
	"time"
)

func newEndpoint(repo Repo.Repo, resources []Resource.Resource) (Endpoint, error) {
	endpoint := Endpoint{
		Repo: repo,
	}

	var err error
	endpoint.Container, err = getContainer(endpoint, resources)
	if err != nil {
		return Endpoint{}, err
	}

	return endpoint, nil
}

func getContainer(e Endpoint, resources []Resource.Resource) (*docker.Container, error) {
	log.WithFields(log.Fields{
		"endpoint": e.Name,
	}).Info("Creating endpoint container")

	container, err := e.SimpleDocker.CreateContainer(
		fmt.Sprintf("%s-endpoint", e.Name),
		fmt.Sprintf("ihsw/%s", e.Name),
		[]string{"DATABASE_HOST=Db"},
	)
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint": e.Name,
			"err":      err.Error(),
		}).Warn("Could not create container")

		return nil, err
	}

	containerLinks := []string{}
	for _, resource := range resources {
		containerLinks = append(
			containerLinks,
			fmt.Sprintf("%s:%s", resource.Container.ID, resource.EndpointTarget),
		)
	}

	if err := e.SimpleDocker.StartContainer(container, containerLinks); err != nil {
		log.WithFields(log.Fields{
			"endpoint": e.Name,
			"err":      err.Error(),
		}).Warn("Could not start container")

		return nil, err
	}

	// waiting for the container to boot up
	time.Sleep(10 * time.Second)

	return container, nil
}

// Endpoint - a container ran against an Endpoint
type Endpoint struct {
	Repo.Repo
	Container *docker.Container
}

// Clean - stops and removes an Endpoint's container
func (e Endpoint) Clean(prevErr error) error {
	if err := e.SimpleDocker.StopContainer(e.Container); err != nil {
		return err
	}

	if err := e.SimpleDocker.RemoveContainer(e.Container); err != nil {
		return err
	}

	return prevErr
}
