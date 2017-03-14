package simpledocker

import (
	"testing"

	docker "github.com/fsouza/go-dockerclient"
)

func TestNewDockerClient(t *testing.T) {
	_, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
	}
}

func TestListImages(t *testing.T) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
	}

	_, err = client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		t.Errorf("Could could not list docker images: %s", err.Error())
	}
}

func TestNewClient(t *testing.T) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
	}

	_ = NewClient(client)
}
