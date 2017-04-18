package resource

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/simpledocker"
)

// DefaultTestResourceName - common resource name for testing
const DefaultTestResourceName = "db"

// CreateTestResource - common test func for creating a resource
func CreateTestResource(client simpledocker.Client, relativePath string, name string, network *docker.Network) (Resource, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return Resource{}, err
	}
	resourceDir, err := filepath.Abs(fmt.Sprintf("%s/%s", cwd, relativePath))
	if err != nil {
		return Resource{}, err
	}

	resource, err := NewResource(client, Opts{
		Name:                 name,
		DockerfileContextDir: resourceDir,
		Network:              network,
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
