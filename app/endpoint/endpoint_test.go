package endpoint

import (
	"fmt"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/galactic-filament/the-matrix/app/repo"
	"github.com/galactic-filament/the-matrix/app/resource"
	"github.com/galactic-filament/the-matrix/app/simpledocker"
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
	endpointResource, err := resource.CreateTestResource(resource.CreateTestResourceOpts{
		Client:       client,
		Network:      endpointNetwork,
		Name:         resource.DefaultTestResourceName,
		RelativePath: fmt.Sprintf("../../%s", resource.DefaultTestResourceName),
	})
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
}
