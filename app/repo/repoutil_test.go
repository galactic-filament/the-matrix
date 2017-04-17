package repo

import (
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/simpledocker"
)

func TestNewRepos(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := simpledocker.NewClient(dockerClient)

	repoNames := []string{DefaultTestRepoName}
	if _, err := NewRepos(repoNames, client); err != nil {
		t.Errorf("Could not create new repos from repo %s: %s", DefaultTestRepoName, err.Error())
	}
}
