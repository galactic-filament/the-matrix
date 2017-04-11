package resource

import (
	"fmt"
	"os"
	"path/filepath"
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
	resourceDir, err := filepath.Abs(fmt.Sprintf("%s/../../%s", cwd, defaultResourceName))
	if err != nil {
		t.Errorf("Could not generate abs resource dir filepath: %s", err.Error())
		return
	}

	resources, err := NewResources(client, []Opts{Opts{
		Name:                 defaultResourceName,
		DockerfileContextDir: resourceDir,
	}})
	if err != nil {
		t.Errorf("Could not create new resources with default resource %s: %s", defaultResourceName, err.Error())
		return
	}

	if err := resources.Clean(); err != nil {
		t.Errorf("Could not clean resources: %s", err.Error())
		return
	}
}
