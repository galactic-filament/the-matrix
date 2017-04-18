package client

import (
	"fmt"
	"testing"

	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/endpoint"
	"github.com/ihsw/the-matrix/app/repo"
	"github.com/ihsw/the-matrix/app/resource"
	"github.com/ihsw/the-matrix/app/simpledocker"
)

func TestRun(t *testing.T) {
	// creating the simpledocker client
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := simpledocker.NewClient(dockerClient)

	// creating the endpoint network
	clientNetwork, err := simpledocker.CreateTestNetwork(client, simpledocker.DefaultTestNetworkName, simpledocker.DefaultTestNetworkDriver)
	if err != nil {
		t.Errorf("Could not create network: %s", err.Error())
		return
	}
	defer simpledocker.CleanupNetwork(t, client, clientNetwork)

	// creating the endpoint resource
	endpointResource, err := resource.CreateTestResource(resource.CreateTestResourceOpts{
		Client:        client,
		Network:       clientNetwork,
		RelativePath:  fmt.Sprintf("../../%s", resource.DefaultTestResourceName),
		Name:          resource.DefaultTestResourceName,
		EndpointLabel: "DATABASE",
	})
	if err != nil {
		t.Errorf("Could not create resource: %s", err.Error())
		return
	}
	defer resource.CleanResource(t, endpointResource)

	// creating the endpoint
	clientEndpoint, err := endpoint.CreateTestEndpoint(endpoint.CreateTestEndpointOpts{
		Client:   client,
		RepoName: repo.DefaultTestRepoName,
		Resource: endpointResource,
		Network:  clientNetwork,
	})
	if err != nil {
		t.Errorf("Could not create endpoint: %s", err.Error())
		return
	}
	defer func() {
		endpoint.CleanEndpoint(t, clientEndpoint)
		time.Sleep(5 * time.Second)
	}()

	// creating a client
	clientRepo, err := repo.NewRepo(DefaultTestClientName, client)
	if err != nil {
		t.Errorf("Could not create repo: %s", err.Error())
		return
	}
	c := NewClient(clientRepo, clientNetwork)

	// running it against the endpoint
	clientContainer, err := c.Run(clientEndpoint)
	if err != nil && err != ErrClientFailed {
		t.Errorf("Could not run endpoint against client: %s", err.Error())
		return
	}
	defer CleanClient(t, c, clientContainer)

	if err == ErrClientFailed {
		containerOutput, err := c.Client.GetContainerLogs(clientContainer)
		if err != nil {
			t.Errorf("Could not get container logs: %s", err.Error())
			return
		}

		t.Errorf("Client failed: %s", containerOutput)
	}
}
