package simpledocker

import (
	"fmt"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	uuid "github.com/nu7hatch/gouuid"
)

const defaultTestContainerName = "test-container"
const defaultTestImage = "hello-world"

func getTestContainerName(prefix string) (string, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s", prefix, u4), nil
}

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

func TestCreateContaienr(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
	}

	client := NewClient(dockerClient)

	testContainerName, err := getTestContainerName(defaultTestContainerName)
	if err != nil {
		t.Errorf("Could not generate test container name: %s", err.Error())
	}
	container, err := client.CreateContainer(testContainerName, defaultTestImage, []string{})
	if err != nil {
		t.Errorf("Could not create %s container: %s", testContainerName, err.Error())
	}

	err = client.RemoveContainer(container)
	if err != nil {
		t.Errorf("Could not remove container %s: %s", testContainerName, err.Error())
	}
}

func TestStartContaienr(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
	}

	client := NewClient(dockerClient)

	testContainerName, err := getTestContainerName(defaultTestContainerName)
	if err != nil {
		t.Errorf("Could not generate test container name: %s", err.Error())
	}
	container, err := client.CreateContainer(testContainerName, defaultTestImage, []string{})
	if err != nil {
		t.Errorf("Could not create %s container: %s", testContainerName, err.Error())
	}

	err = client.StartContainer(container, []string{})
	if err != nil {
		t.Errorf("Could not start container %s: %s", testContainerName, err.Error())
	}

	err = client.RemoveContainer(container)
	if err != nil {
		t.Errorf("Could not remove container %s: %s", testContainerName, err.Error())
	}
}
