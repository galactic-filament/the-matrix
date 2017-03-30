package simpledocker

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/util"
)

const defaultTestImageOutput = "hello world\n"
const defaultTestContainerName = "test-container"
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

func TestCreateContainer(t *testing.T) {
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
	defer cleanupContainer(t, client, container)
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
	defer cleanupContainer(t, client, container)

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

	containerName, container, err := createTestContainer(client, defaultTestContainerName, defaultTestImage, []string{})
	if err != nil {
		t.Errorf("Could not create %s container: %s", containerName, err.Error())
		return
	}
	defer cleanupContainer(t, client, container)

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

func TestHasImage(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	nonexistentImageName, err := util.GetPrefixedUUID("fdsfgs")
	if err != nil {
		t.Errorf("Could not create non-existent image name fdsfgs: %s", err.Error())
		return
	}

	hasImage, err := client.HasImage(nonexistentImageName)
	if err != nil {
		t.Errorf("Could not check for image %s: %s", nonexistentImageName, err.Error())
		return
	}

	if hasImage {
		t.Errorf("Image %s was found when it should not have been", nonexistentImageName)
		return
	}
}

func TestPullImage(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	if err := client.PullImage(defaultTestImage, defaultTestImageTag); err != nil {
		t.Errorf("Could not pull test image %s: %s", defaultTestImage, err.Error())
	}
	cleanupImage(t, client, defaultTestImage)
}

func TestRemoveImage(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	hasImage, err := client.HasImage(defaultTestImage)
	if err != nil {
		t.Errorf("Could not check if image exists %s: %s", defaultTestImage, err.Error())
		return
	}

	if !hasImage {
		if err := client.PullImage(defaultTestImage, defaultTestImageTag); err != nil {
			t.Errorf("Could not pull default test image %s with tag %s: %s", defaultTestImage, defaultTestImageTag, err.Error())
			return
		}
	}

	imageID := fmt.Sprintf("%s:%s", defaultTestImage, defaultTestImageTag)
	if err := client.RemoveImage(imageID); err != nil {
		t.Errorf("Could not remove image %s: %s", imageID, err.Error())
		return
	}
}

func TestBuildImage(t *testing.T) {
	// creating a new simpledocker client
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
	cleanupImage(t, client, exampleImageName)
}

func TestIsRunning(t *testing.T) {
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
	defer cleanupImage(t, client, defaultTestImage)

	// creating one and starting it up
	name, err := util.GetPrefixedUUID(defaultTestContainerName)
	if err != nil {
		t.Errorf("Could not create prefixed container name: %s", err.Error())
		return
	}
	container, err := client.CreateContainer(name, defaultTestImage, []string{})
	if err != nil {
		t.Errorf("Could not create test container from image %s: %s", defaultTestImage, err.Error())
		return
	}
	defer cleanupContainer(t, client, container)

	if err := client.StartContainer(container, []string{}); err != nil {
		t.Errorf("Could not start container %s: %s", container.Name, err.Error())
		return
	}

	// waiting 5s for it to exit
	time.Sleep(5 * time.Second)

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
	defer cleanupImage(t, client, defaultDbImage)

	// creating one and starting it up
	name, err := util.GetPrefixedUUID(defaultTestContainerName)
	if err != nil {
		t.Errorf("Could not create prefixed container name: %s", err.Error())
		return
	}
	container, err := client.CreateContainer(name, defaultDbImage, []string{})
	if err != nil {
		t.Errorf("Could not create test container from image %s: %s", defaultDbImage, err.Error())
		return
	}
	defer cleanupContainer(t, client, container)

	if err := client.StartContainer(container, []string{}); err != nil {
		t.Errorf("Could not start container %s: %s", container.Name, err.Error())
		return
	}

	// waiting 5s to see if it is still up
	time.Sleep(5 * time.Second)

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
	defer cleanupImage(t, client, exampleImageName)

	// creating the container
	name, err := util.GetPrefixedUUID(defaultTestContainerName)
	if err != nil {
		t.Errorf("Could not create prefixed container name: %s", err.Error())
		return
	}
	container, err := client.CreateContainer(name, exampleImageName, []string{})
	if err != nil {
		t.Errorf("Could not create a container from the test image: %s", err.Error())
		return
	}
	defer cleanupContainer(t, client, container)

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
