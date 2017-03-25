package resource

import (
	"fmt"
	"os"
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

	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Could not get working dir: %s", err.Error())
		return
	}
	resourceDir := fmt.Sprintf("%s/../../%s", cwd, defaultResourceName)

	resources, err := NewResources(client, []Opts{Opts{
		Name:                 defaultResourceName,
		DockerfileContextDir: resourceDir,
		EndpointEnvVars:      map[string]string{"DATABASE_HOST": "Db"}},
	})
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
