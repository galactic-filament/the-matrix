package simpledocker

import (
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/util"
)

const defaultTestContainerName = "test-container"
const defaultTestImageOutput = "hello world\n"
const defaultTestImage = "hello-world"
const defaultTestImageTag = "latest"
const defaultDbImage = "postgres"

func createTestContainer(client Client, namePrefix string, imageName string, links []string) (string, *docker.Container, error) {
	containerName, err := util.GetPrefixedUUID(namePrefix)
	if err != nil {
		return "", nil, err
	}

	hasImage, err := client.HasImage(defaultTestImage)
	if err != nil {
		return "", nil, err
	}

	if !hasImage {
		if err := client.PullImage(defaultTestImage, defaultTestImageTag); err != nil {
			return "", nil, err
		}
	}

	container, err := client.CreateContainer(containerName, imageName, links)
	return containerName, container, err
}

func cleanupContainer(t *testing.T, client Client, container *docker.Container) {
	isRunning, err := client.IsRunning(container)
	if err != nil {
		t.Errorf("Could not check if container was running: %s", err.Error())
		return
	}
	if isRunning {
		if err := client.StopContainer(container); err != nil {
			t.Errorf("Could not stop container: %s", err.Error())
			return
		}
	}

	if err := client.RemoveContainer(container); err != nil {
		t.Errorf("Could not remove container: %s", err.Error())
		return
	}
}

func cleanupImage(t *testing.T, client Client, imageID string) {
	if err := client.RemoveImage(imageID); err != nil {
		t.Errorf("Could not remove image: %s", err.Error())
		return
	}
}

func cleanupNetwork(t *testing.T, client Client, network *docker.Network) {
	if err := client.RemoveNetwork(network); err != nil {
		t.Errorf("Could not remove network: %s", err.Error())
		return
	}
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
		return
	}

	_, err = client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		t.Errorf("Could could not list docker images: %s", err.Error())
		return
	}
}
