package resource

import (
	"errors"
	"fmt"

	"github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/simpledocker"
)

func getContainerID(name string) string      { return fmt.Sprintf("%s-resource", name) }
func getContainerImageID(name string) string { return fmt.Sprintf("ihsw/the-matrix-%s", name) }

func newResource(client simpledocker.Client, name string) (Resource, error) {
	r := Resource{client, name, nil}

	container, err := client.CreateContainer(
		getContainerID(name),
		getContainerImageID(name),
		[]string{},
	)
	if err != nil {
		return Resource{}, err
	}
	r.container = container

	if err := r.client.StartContainer(container, []string{}); err != nil {
		return Resource{}, err
	}

	return r, nil
}

// Resource - a container for each Endpoint to use (database, etc)
type Resource struct {
	client    simpledocker.Client
	name      string
	container *docker.Container
}

// Clean - stops and removes the Resource's container
func (r Resource) Clean() error {
	if r.container == nil {
		return errors.New("Resource container was nil")
	}

	if err := r.client.StopContainer(r.container); err != nil {
		return err
	}

	if err := r.client.RemoveContainer(r.container); err != nil {
		return err
	}

	return nil
}

// GetLinkLine - returns the expected docker link line
func (r Resource) GetLinkLine() string { return fmt.Sprintf("%s:%s", r.name, r.name) }
