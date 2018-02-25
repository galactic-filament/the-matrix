package resource

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/galactic-filament/the-matrix/app/simpledocker"
)

// DefaultTestResourceName - common resource name for testing
const DefaultTestResourceName = "db"

// CreateTestResourceOpts - opts for corresponding func
type CreateTestResourceOpts struct {
	Client        simpledocker.Client
	Network       *docker.Network
	RelativePath  string
	Name          string
	EndpointLabel string
}

// CreateTestResource - common test func for creating a resource
func CreateTestResource(opts CreateTestResourceOpts) (Resource, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return Resource{}, err
	}
	resourceDir, err := filepath.Abs(fmt.Sprintf("%s/%s", cwd, opts.RelativePath))
	if err != nil {
		return Resource{}, err
	}

	resource, err := NewResource(opts.Client, Opts{
		Name:                 opts.Name,
		DockerfileContextDir: resourceDir,
		Network:              opts.Network,
		EndpointLabel:        opts.EndpointLabel,
	})
	if err != nil {
		return Resource{}, err
	}

	return resource, nil
}

// CleanResource - common test func used for cleaning up a resource
func CleanResource(t *testing.T, resource Resource) {
	if err := resource.Clean(); err != nil {
		t.Errorf("Could not clean resource %s: %s", resource.name, err.Error())
		return
	}
}

// CleanResources - common test func used for cleaning up resources
func CleanResources(t *testing.T, resources Resources) {
	if err := resources.Clean(); err != nil {
		t.Errorf("Could not clean resources: %s", err.Error())
		return
	}
}
