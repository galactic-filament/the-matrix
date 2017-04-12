package resource

import (
	"path/filepath"
	"testing"

	"fmt"
	"os"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/simpledocker"
)

const defaultResourceName = "db"
const defaultTestNetworkName = "galaxy"
const defaultTestNetworkDriver = "bridge"

func getTestResource(client simpledocker.Client, relativePath string, name string, network *docker.Network) (Resource, error) {
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
	})
	if err != nil {
		return Resource{}, err
	}

	return resource, nil
}

func cleanResource(t *testing.T, resource Resource) {
	if err := resource.Clean(); err != nil {
		t.Errorf("Could not clean resource %s: %s", resource.name, err.Error())
		return
	}
}

func TestNewResource(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := simpledocker.NewClient(dockerClient)

	resource, err := getTestResource(client, fmt.Sprintf("../../%s", defaultResourceName), defaultResourceName, nil)
	if err != nil {
		t.Errorf("Could not create test resource: %s", err.Error())
		return
	}
	defer cleanResource(t, resource)
}

func TestGetContainerIP(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := simpledocker.NewClient(dockerClient)

	network, err := client.CreateNetwork(defaultTestNetworkName, defaultTestNetworkDriver)
	if err != nil {
		t.Errorf("Could not create network: %s", err.Error())
		return
	}
	defer simpledocker.CleanupNetwork(t, client, network)

	resource, err := getTestResource(client, fmt.Sprintf("../../%s", defaultResourceName), defaultResourceName, network)
	if err != nil {
		t.Errorf("Could not create test resource: %s", err.Error())
		return
	}
	defer cleanResource(t, resource)

	ip, err := resource.GetContainerIP()
	if err != nil {
		t.Errorf("Could not get resource IP: %s", err.Error())
		return
	}
	if ip == nil {
		t.Errorf("Resource IP was nil")
		return
	}
}
