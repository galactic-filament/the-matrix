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

func createTestContainer(client Client, namePrefix string, imageName string, links []string) (string, *docker.Container, error) {
	containerName, err := getTestContainerName(namePrefix)
	if err != nil {
		return "", nil, err
	}

	container, err := client.CreateContainer(containerName, imageName, links)
	return containerName, container, err
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

func TestCreateContaienr(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	containerName, container, err := createTestContainer(client, defaultTestContainerName, defaultTestImage, []string{})
	if err != nil {
		t.Errorf("Could not create %s container: %s", containerName, err.Error())
		return
	}

	err = client.RemoveContainer(container)
	if err != nil {
		t.Errorf("Could not remove container %s: %s", container.Name, err.Error())
		return
	}
}

func TestStartContainer(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	containerName, container, err := createTestContainer(client, defaultTestContainerName, defaultTestImage, []string{})
	if err != nil {
		t.Errorf("Could not create %s container: %s", containerName, err.Error())
		return
	}

	err = client.StartContainer(container, []string{})
	if err != nil {
		t.Errorf("Could not start container %s: %s", container.Name, err.Error())
		return
	}

	_, err = dockerClient.WaitContainer(container.ID)
	if err != nil {
		t.Errorf("Could not wait to exit container %s: %s", container.Name, err.Error())
		return
	}

	err = client.RemoveContainer(container)
	if err != nil {
		t.Errorf("Could not remove container %s: %s", container.Name, err.Error())
		return
	}
}

func TestRunContainer(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	containerName, container, err := createTestContainer(client, defaultTestContainerName, defaultTestImage, []string{})
	if err != nil {
		t.Errorf("Could not create %s container: %s", containerName, err.Error())
		return
	}

	// starting it up via run
	type runContainerResult struct {
		failed bool
		err    error
	}
	runOut := make(chan runContainerResult)
	go func() {
		failed, err := client.RunContainer(container, []string{})
		runOut <- runContainerResult{failed, err}
	}()

	// gathering the result
	result := <-runOut
	if err := result.err; err != nil {
		t.Errorf("Could not run container %s: %s", container.Name, err.Error())
		return
	}

	// cleaning it up
	err = client.RemoveContainer(container)
	if err != nil {
		t.Errorf("Could not remove container %s: %s", container.Name, err.Error())
		return
	}
}
