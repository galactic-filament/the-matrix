package repo

import (
	"testing"

	"github.com/Sirupsen/logrus"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/simpledocker"
)

func TestNewRepo(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := simpledocker.NewClient(dockerClient)

	if _, err := NewRepo(DefaultTestRepoName, client); err != nil {
		t.Errorf("Could not create new repo %s: %s", DefaultTestRepoName, err.Error())
	}
}

func TestNonexistentNewRepo(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := simpledocker.NewClient(dockerClient)

	imageName := GetImageName(DefaultTestRepoName)
	hasImage, err := client.HasImage(imageName)
	if err != nil {
		t.Errorf("Could not validate whether default repo image exists %s: %s", imageName, err.Error())
		return
	}

	if hasImage {
		t.Logf("Removing image %s to be pulled by newRepo()", imageName)

		containers, err := client.GetContainersByImageID(imageName, defaultRepoImageTag)
		if err != nil {
			t.Errorf("Could not get containers for this image: %s", err.Error())
			return
		}
		logrus.WithFields(logrus.Fields{"containers": len(containers), "image": imageName}).Info("Containers found with image")
		for _, container := range containers {
			logrus.WithFields(logrus.Fields{"container": container.Name}).Info("Container found")
		}

		if err := client.RemoveImage(imageName); err != nil {
			t.Errorf("Could not remove default repo image %s :%s", imageName, err.Error())
			return
		}
	}

	t.Logf("Creating new repo %s", DefaultTestRepoName)
	if _, err := NewRepo(DefaultTestRepoName, client); err != nil {
		t.Errorf("Could not create new repo %s :%s", DefaultTestRepoName, err.Error())
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

	imageName := GetImageName(DefaultTestRepoName)
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

	if _, err := NewRepo(DefaultTestRepoName, client); err != nil {
		t.Errorf("Could not create new repo %s :%s", DefaultTestRepoName, err.Error())
		return
	}
}
