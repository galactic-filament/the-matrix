package resource

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/simpledocker"
	"github.com/ihsw/the-matrix/app/util"
)

func getContainerID(name string) string { return fmt.Sprintf("%s-resource", name) }
func getImageName(name string) string   { return fmt.Sprintf("ihsw/the-matrix-%s", name) }

// NewResource - creates a new resource based on a dockerfile, optionally building it where it does not exist
func NewResource(client simpledocker.Client, opts Opts) (Resource, error) {
	r := Resource{
		client:        client,
		network:       opts.Network,
		name:          opts.Name,
		EndpointLabel: opts.EndpointLabel,
	}
	imageID := getImageName(r.name)

	// validating that the resource image exists, if not then building it
	hasImage, err := client.HasImage(imageID)
	if err != nil {
		return Resource{}, err
	}
	if !hasImage {
		if err := client.BuildImage(imageID, opts.DockerfileContextDir); err != nil {
			return Resource{}, err
		}
	}

	// creating an endpoint container
	containerID, err := util.GetPrefixedUUID(getContainerID(r.name))
	if err != nil {
		return Resource{}, err
	}
	container, err := client.CreateContainer(simpledocker.CreateContainerOptions{
		Name:    containerID,
		Image:   imageID,
		Network: opts.Network,
	})
	if err != nil {
		return Resource{}, err
	}
	r.container = container

	// starting it up
	if err := r.client.StartContainer(container, []string{}); err != nil {
		return Resource{}, err
	}

	// waiting for the resource to start up
	time.Sleep(5 * time.Second)

	return r, nil
}

// Opts - expected values when calling new-resource
type Opts struct {
	Name                 string
	DockerfileContextDir string
	Network              *docker.Network
	EndpointLabel        string
}

// Resource - a container for each Endpoint to use (database, etc)
type Resource struct {
	client        simpledocker.Client
	network       *docker.Network
	container     *docker.Container
	name          string
	EndpointLabel string
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

// GetContainerIP - returns ip address of resource container in the resource's network
func (r Resource) GetContainerIP() (net.IP, error) {
	return r.client.GetContainerIP(r.network, r.container)
}
