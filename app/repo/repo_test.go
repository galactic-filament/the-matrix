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

	if _, err := NewRepo(defaultRepoName, client); err != nil {
		t.Errorf("Could not create new repo %s: %s", defaultRepoName, err.Error())
	}
}

func TestNonexistentNewRepo(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := simpledocker.NewClient(dockerClient)

	imageName := GetImageName(defaultRepoName)
	hasImage, err := client.HasImage(imageName)
	if err != nil {
		t.Errorf("Could not validate whether default repo image exists %s: %s", imageName, err.Error())
		return
	}

	if hasImage {
		t.Logf("Removing image %s to be pulled by newRepo()", imageName)
		if err := client.RemoveImage(imageName); err != nil {
			t.Errorf("Could not remove default repo image %s :%s", imageName, err.Error())
			return
		}
	}

	t.Logf("Creating new repo %s", defaultRepoName)
	if _, err := NewRepo(defaultRepoName, client); err != nil {
		t.Errorf("Could not create new repo %s :%s", defaultRepoName, err.Error())
		return
	}
}

func TestExistentNewRepo(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := simpledocker.NewClient(dockerClient)

	imageName := GetImageName(defaultRepoName)
	hasImage, err := client.HasImage(imageName)
	if err != nil {
		t.Errorf("Could not validate whether default repo image exists %s: %s", imageName, err.Error())
		return
	}

	if !hasImage {
		t.Logf("Repo image %s was not found, pulling manually", imageName)
		err = client.PullImage(imageName, defaultRepoImageTag)
		if err != nil {
			t.Errorf("Could not pull default repo image %s: %s", imageName, err.Error())
			return
		}
	}

	if _, err := NewRepo(defaultRepoName, client); err != nil {
		t.Errorf("Could not create new repo %s :%s", defaultRepoName, err.Error())
		return
	}
}
