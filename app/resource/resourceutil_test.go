package resource

import (
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/simpledocker"
)

func TestNewResources(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := simpledocker.NewClient(dockerClient)

	resources, err := NewResources(client, []string{defaultResourceName})
	if err != nil {
		t.Errorf("Could not create new resources with default resource %s: %s", defaultResourceName, err.Error())
		return
	}

	for _, resource := range resources {
		if err := resource.Clean(); err != nil {
			t.Errorf("Could not clean resource %s: %s", resource.name, err.Error())
			return
		}
	}
}
