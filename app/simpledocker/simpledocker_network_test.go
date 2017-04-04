package simpledocker

import (
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/util"
)

func TestCreateNetwork(t *testing.T) {
	// creating a simpledocker client
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	networkName, err := util.GetPrefixedUUID("test-network")
	if err != nil {
		t.Errorf("Could not generate network name: %s", err.Error())
		return
	}
	network, err := client.CreateNetwork(networkName, defaultNetworkDriver)
	if err != nil {
		t.Errorf("Could not create network: %s", err.Error())
		return
	}
	defer cleanupNetwork(t, client, network)
}
