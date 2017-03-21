package repo

import (
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/simpledocker"
)

const defaultRepoName = "es-bueno"

func TestNewRepo(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := simpledocker.NewClient(dockerClient)

	if _, err := newRepo(defaultRepoName, client); err != nil {
		t.Errorf("Could not create new repo %s: %s", defaultRepoName, err.Error())
	}
}
