package client

import (
	"fmt"
	"testing"

	"github.com/Sirupsen/logrus"
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
	logrus.Info("Creating test network")
	clientNetwork, err := simpledocker.CreateTestNetwork(client, simpledocker.DefaultTestNetworkName, simpledocker.DefaultTestNetworkDriver)
	if err != nil {
		t.Errorf("Could not create network: %s", err.Error())
		return
	}
	defer simpledocker.CleanupNetwork(t, client, clientNetwork)

	// creating the endpoint resource
	logrus.Info("Creating test resource")
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
	logrus.Info("Creating test endpoint")
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
	defer endpoint.CleanEndpoint(t, clientEndpoint)

	// creating a client
	logrus.Info("Creating a test client")
	clientRepo, err := repo.NewRepo(DefaultTestClientName, client)
	if err != nil {
		t.Errorf("Could not create repo: %s", err.Error())
		return
	}
	c := NewClient(clientRepo, clientNetwork)

	// running it against the endpoint
	logrus.Info("Running it against the endpoint")
	clientContainer, err := c.Run(clientEndpoint)
	if err != nil && err != ErrClientFailed {
		t.Errorf("Could not run endpoint against client: %s", err.Error())
		return
	}
	defer CleanClient(t, c, clientContainer)

	logrus.Info("Checking for failure")
	if err == ErrClientFailed {
		containerOutput, err := c.Client.GetContainerLogs(clientContainer)
		if err != nil {
			t.Errorf("Could not get container logs: %s", err.Error())
			return
		}

		t.Errorf("Client failed: %s", containerOutput)
	}
}
