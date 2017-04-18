package resource

import (
	"fmt"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/simpledocker"
)

func TestNewResource(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := simpledocker.NewClient(dockerClient)

	resource, err := CreateTestResource(client, fmt.Sprintf("../../%s", DefaultTestResourceName), DefaultTestResourceName, nil)
	if err != nil {
		t.Errorf("Could not create test resource: %s", err.Error())
		return
	}
	defer CleanResource(t, resource)
}

func TestGetContainerIP(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := simpledocker.NewClient(dockerClient)

	network, err := simpledocker.CreateTestNetwork(client, simpledocker.DefaultTestNetworkName, simpledocker.DefaultTestNetworkDriver)
	if err != nil {
		t.Errorf("Could not create network: %s", err.Error())
		return
	}
	defer simpledocker.CleanupNetwork(t, client, network)

	resource, err := CreateTestResource(client, fmt.Sprintf("../../%s", DefaultTestResourceName), DefaultTestResourceName, network)
	if err != nil {
		t.Errorf("Could not create test resource: %s", err.Error())
		return
	}
	defer CleanResource(t, resource)

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
