package resource

import (
	"errors"
	"fmt"

	"github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/simpledocker"
)

func getContainerID(name string) string { return fmt.Sprintf("%s-resource", name) }
func getImageID(name string) string     { return fmt.Sprintf("ihsw/the-matrix-%s", name) }

func newResource(client simpledocker.Client, opts Opts) (Resource, error) {
	r := Resource{client, opts.Name, nil}
	imageID := getImageID(r.name)

	hasImage, err := client.HasImage(imageID)
	if err != nil {
		return Resource{}, err
	}
	if !hasImage {
		if err := client.BuildImage(imageID, opts.DockerfileContextDir); err != nil {
			return Resource{}, err
		}
	}

	container, err := client.CreateContainer(
		getContainerID(r.name),
		imageID,
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

// Opts - expect values when calling new-resource
type Opts struct {
	Name                 string
	DockerfileContextDir string
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
