package simpledocker

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/util"
)

const defaultTestContainerName = "test-container"

func createTestContainer(client Client, namePrefix string, imageName string, network *docker.Network) (string, *docker.Container, error) {
	containerName, err := util.GetPrefixedUUID(namePrefix)
	if err != nil {
		return "", nil, err
	}

	hasImage, err := client.HasImage(imageName)
	if err != nil {
		return "", nil, err
	}

	if !hasImage {
		if err := client.PullImage(imageName, defaultTestImageTag); err != nil {
			return "", nil, err
		}
	}

	container, err := client.CreateContainer(CreateContainerOptions{Name: containerName, Image: imageName, Network: network})
	return containerName, container, err
}

func TestCreateContainer(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	containerName, container, err := createTestContainer(client, defaultTestContainerName, defaultTestImage, nil)
	if err != nil {
		t.Errorf("Could not create %s container: %s", containerName, err.Error())
		return
	}
	defer CleanupContainer(t, client, container)
}

func TestStartContainer(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	containerName, container, err := createTestContainer(client, defaultTestContainerName, defaultTestImage, nil)
	if err != nil {
		t.Errorf("Could not create %s container: %s", containerName, err.Error())
		return
	}
	defer CleanupContainer(t, client, container)

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
}

func TestRunContainer(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	containerName, container, err := createTestContainer(client, defaultTestContainerName, defaultTestImage, nil)
	if err != nil {
		t.Errorf("Could not create %s container: %s", containerName, err.Error())
		return
	}
	defer CleanupContainer(t, client, container)

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
}

func TestIsNotRunning(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	// pulling the default test image where appropriate
	hasImage, err := client.HasImage(defaultTestImage)
	if err != nil {
		t.Errorf("Could not check if has image: %s", err.Error())
		return
	}
	if !hasImage {
		if err := client.PullImage(defaultTestImage, defaultTestImageTag); err != nil {
			t.Errorf("Could not pull image %s: %s", defaultTestImage, err.Error())
			return
		}
	}
	defer CleanupImage(t, client, defaultTestImage)

	// creating one and starting it up
	name, err := util.GetPrefixedUUID(defaultTestContainerName)
	if err != nil {
		t.Errorf("Could not create prefixed container name: %s", err.Error())
		return
	}
	container, err := client.CreateContainer(CreateContainerOptions{Name: name, Image: defaultTestImage})
	if err != nil {
		t.Errorf("Could not create test container from image %s: %s", defaultTestImage, err.Error())
		return
	}
	defer CleanupContainer(t, client, container)
	if err := client.StartContainer(container, []string{}); err != nil {
		t.Errorf("Could not start container %s: %s", container.Name, err.Error())
		return
	}

	// waiting for it to exit
	if _, err := client.dockerClient.WaitContainer(container.ID); err != nil {
		t.Errorf("Could not wait for container to exit: %s", err.Error())
		return
	}

	// validating that it exited
	isRunning, err := client.IsRunning(container)
	if err != nil {
		t.Errorf("Could not check if container %s is running: %s", container.Name, err.Error())
		return
	}
	if isRunning {
		t.Errorf("Container %s was found to be running when it should not have been", container.Name)
		return
	}
}

func TestIsStillRunning(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	// pulling the default test image where appropriate
	hasImage, err := client.HasImage(defaultDbImage)
	if err != nil {
		t.Errorf("Could not check if has image: %s", err.Error())
		return
	}
	if !hasImage {
		if err := client.PullImage(defaultDbImage, defaultTestImageTag); err != nil {
			t.Errorf("Could not pull image %s: %s", defaultDbImage, err.Error())
			return
		}
	}

	// creating one and starting it up
	name, err := util.GetPrefixedUUID(defaultTestContainerName)
	if err != nil {
		t.Errorf("Could not create prefixed container name: %s", err.Error())
		return
	}
	container, err := client.CreateContainer(CreateContainerOptions{Name: name, Image: defaultDbImage})
	if err != nil {
		t.Errorf("Could not create test container from image %s: %s", defaultDbImage, err.Error())
		return
	}
	defer CleanupContainer(t, client, container)

	if err := client.StartContainer(container, []string{}); err != nil {
		t.Errorf("Could not start container %s: %s", container.Name, err.Error())
		return
	}

	// waiting 5s to see if it is still up
	// time.Sleep(5 * time.Second)

	// validating that it is still up
	isRunning, err := client.IsRunning(container)
	if err != nil {
		t.Errorf("Could not check if container %s is running: %s", container.Name, err.Error())
		return
	}
	if !isRunning {
		t.Errorf("Container %s was found to be not running when it should have been", container.Name)
		return
	}
}

func TestGetContainerLogs(t *testing.T) {
	// creating a simpledocker client
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	// gathering context dir for the dockerfile
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Could not get working dir: %s", err.Error())
		return
	}
	contextDir, err := filepath.Abs(fmt.Sprintf("%s/../../test-fixtures", cwd))
	if err != nil {
		t.Errorf("Could not get absolute filepath for test fixture: %s", err.Error())
		return
	}

	// generating an image name
	exampleImageName, err := util.GetPrefixedUUID("hello-world")
	if err != nil {
		t.Errorf("Could not generate example image name: %s", err.Error())
		return
	}

	// building the image
	if err := client.BuildImage(exampleImageName, contextDir); err != nil {
		t.Errorf("Could not build example image: %s", err.Error())
		return
	}
	defer CleanupImage(t, client, exampleImageName)

	// creating the container
	name, err := util.GetPrefixedUUID(defaultTestContainerName)
	if err != nil {
		t.Errorf("Could not create prefixed container name: %s", err.Error())
		return
	}
	container, err := client.CreateContainer(CreateContainerOptions{Name: name, Image: exampleImageName})
	if err != nil {
		t.Errorf("Could not create a container from the test image: %s", err.Error())
		return
	}
	defer CleanupContainer(t, client, container)

	// starting it up
	if err := client.StartContainer(container, []string{}); err != nil {
		t.Errorf("Could not start test container: %s", err.Error())
		return
	}

	// waiting for it to exit
	if _, err := client.dockerClient.WaitContainer(container.ID); err != nil {
		t.Errorf("Could not wait for container to exit: %s", err.Error())
		return
	}

	// gathering the log output
	containerOutput, err := client.GetContainerLogs(container)
	if err != nil {
		t.Errorf("Could not get container logs: %s", err.Error())
		return
	}
	if containerOutput != defaultTestImageOutput {
		t.Errorf("Container output did not match the expected output: %s vs %s", defaultTestImageOutput, containerOutput)
		return
	}
}
