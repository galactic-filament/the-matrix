package resource

import (
	"errors"
	"fmt"

	"github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/simpledocker"
	"github.com/ihsw/the-matrix/app/util"
)

func getContainerID(name string) string { return fmt.Sprintf("%s-resource", name) }
func getImageID(name string) string     { return fmt.Sprintf("ihsw/the-matrix-%s", name) }

// NewResource - creates a new resource based on a dockerfile, optionally building it where it does not exist
func NewResource(client simpledocker.Client, opts Opts) (Resource, error) {
	r := Resource{client, opts.Name, nil, opts.EndpointEnvVars}
	imageID := getImageID(r.name)

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
	container, err := client.CreateContainer(containerID, imageID, []string{})
	if err != nil {
		return Resource{}, err
	}
	r.container = container

	// starting it up
	if err := r.client.StartContainer(container, []string{}); err != nil {
		return Resource{}, err
	}

	return r, nil
}

// Opts - expected values when calling new-resource
type Opts struct {
	Name                 string
	DockerfileContextDir string
	EndpointEnvVars      map[string]string
}

// Resource - a container for each Endpoint to use (database, etc)
type Resource struct {
	client          simpledocker.Client
	name            string
	container       *docker.Container
	endpointEnvVars map[string]string
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

// GetEnvVars - returns the env vars map as an array of strings
func (r Resource) GetEnvVars() []string {
	envVars := []string{}
	for k, v := range r.endpointEnvVars {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}

	return envVars
}
