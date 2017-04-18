package endpoint

import (
	"fmt"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/repo"
	"github.com/ihsw/the-matrix/app/resource"
	"github.com/ihsw/the-matrix/app/simpledocker"
)

func TestNewEndpoint(t *testing.T) {
	// creating the simpledocker client
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := simpledocker.NewClient(dockerClient)

	// creating the endpoint network
	endpointNetwork, err := simpledocker.CreateTestNetwork(client, simpledocker.DefaultTestNetworkName, simpledocker.DefaultTestNetworkDriver)
	if err != nil {
		t.Errorf("Could not create network: %s", err.Error())
		return
	}
	defer simpledocker.CleanupNetwork(t, client, endpointNetwork)

	// creating the endpoint resource
	endpointResource, err := resource.CreateTestResource(
		client,
		fmt.Sprintf("../../%s", resource.DefaultTestResourceName),
		resource.DefaultTestResourceName,
		endpointNetwork,
	)
	if err != nil {
		t.Errorf("Could not create resource: %s", err.Error())
		return
	}
	defer resource.CleanResource(t, endpointResource)

	// creating the endpoint
	e, err := CreateTestEndpoint(CreateTestEndpointOpts{
		Client:   client,
		RepoName: repo.DefaultTestRepoName,
		Resource: endpointResource,
		Network:  endpointNetwork,
	})
	if err != nil {
		t.Errorf("Could not create endpoint: %s", err.Error())
		return
	}
	defer CleanEndpoint(t, e)

	// verifying that it is running
	isRunning, err := client.IsRunning(e.Container)
	if err != nil {
		t.Errorf("Could not check if endpoint container is running: %s", err.Error())
		return
	}
	if !isRunning {
		containerOutput, err := client.GetContainerLogs(e.Container)
		if err != nil {
			t.Errorf("Could not fetch container logs: %s", err.Error())
			return
		}

		t.Errorf("Endpoint container %s was not up: %s", e.Name, containerOutput)
		return
	}
}
