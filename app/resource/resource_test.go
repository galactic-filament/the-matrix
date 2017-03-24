package resource

import (
	"testing"

	"fmt"
	"os"

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

	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Could not get working dir: %s", err.Error())
		return
	}
	resourceDir := fmt.Sprintf("%s/../../%s", cwd, defaultResourceName)

	resource, err := newResource(client, Opts{defaultResourceName, resourceDir})
	if err != nil {
		t.Errorf("Could not create repo %s: %s", defaultResourceName, err.Error())
		return
	}

	if err := resource.Clean(); err != nil {
		t.Errorf("Could not clean resource %s: %s", resource.name, err.Error())
		return
	}
}
