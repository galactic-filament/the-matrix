package resource

import (
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/simpledocker"
)

const defaultResourceName = "db"

func TestNewResource(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := simpledocker.NewClient(dockerClient)

	resource, err := newResource(client, defaultResourceName)
	if err != nil {
		t.Errorf("Could not create repo %s: %s", defaultResourceName, err.Error())
		return
	}

	if err := resource.Clean(); err != nil {
		t.Errorf("Could not clean resource %s: %s", resource.name, err.Error())
		return
	}
}
