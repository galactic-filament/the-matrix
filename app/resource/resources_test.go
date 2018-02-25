package resource

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/galactic-filament/the-matrix/app/simpledocker"
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
	resourceDir, err := filepath.Abs(fmt.Sprintf("%s/../../%s", cwd, DefaultTestResourceName))
	if err != nil {
		t.Errorf("Could not generate abs resource dir filepath: %s", err.Error())
		return
	}

	resources, err := NewResources(client, []Opts{Opts{
		Name:                 DefaultTestResourceName,
		DockerfileContextDir: resourceDir,
	}})
	if err != nil {
		t.Errorf("Could not create new resources with default resource %s: %s", DefaultTestResourceName, err.Error())
		return
	}
	defer CleanResources(t, resources)
}
