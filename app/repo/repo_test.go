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

func TestNonexistentNewRepo(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := simpledocker.NewClient(dockerClient)

	imageID := GetImageID(defaultRepoName)
	hasImage, err := client.HasImage(imageID)
	if err != nil {
		t.Errorf("Could not validate whether default repo image exists %s: %s", imageID, err.Error())
		return
	}

	if hasImage {
		t.Logf("Removing image %s to be pulled by newRepo()", imageID)
		if err := client.RemoveImage(imageID); err != nil {
			t.Errorf("Could not remove default repo image %s :%s", imageID, err.Error())
			return
		}
	}

	t.Logf("Creating new repo %s", defaultRepoName)
	if _, err := newRepo(defaultRepoName, client); err != nil {
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

	imageID := GetImageID(defaultRepoName)
	hasImage, err := client.HasImage(imageID)
	if err != nil {
		t.Errorf("Could not validate whether default repo image exists %s: %s", imageID, err.Error())
		return
	}

	if !hasImage {
		t.Logf("Repo image %s was not found, pulling manually", imageID)
		err = client.PullImage(imageID, defaultRepoImageTag)
		if err != nil {
			t.Errorf("Could not pull default repo image %s: %s", imageID, err.Error())
			return
		}
	}

	if _, err := newRepo(defaultRepoName, client); err != nil {
		t.Errorf("Could not create new repo %s :%s", defaultRepoName, err.Error())
		return
	}
}
