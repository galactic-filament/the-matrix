package simpledocker

import (
	"fmt"
	"path/filepath"
	"testing"

	"os"

	"time"

	docker "github.com/fsouza/go-dockerclient"
	uuid "github.com/nu7hatch/gouuid"
)

const defaultTestContainerName = "test-container"
const defaultTestImage = "hello-world"
const defaultTestImageTag = "latest"
const defaultDbImage = "postgres"

func getPrefixedUUID(prefix string) (string, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s", prefix, u4), nil
}

func createTestContainer(client Client, namePrefix string, imageName string, links []string) (string, *docker.Container, error) {
	containerName, err := getPrefixedUUID(namePrefix)
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

func TestHasImage(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	nonexistentImageName, err := getPrefixedUUID("fdsfgs")
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
	exampleImageName, err := getPrefixedUUID("hello-world")
	if err != nil {
		t.Errorf("Could not generate example image name: %s", err.Error())
		return
	}

	// building the image
	if err := client.BuildImage(exampleImageName, contextDir); err != nil {
		t.Errorf("Could not build example image: %s", err.Error())
		return
	}

	// removing the image
	if err := client.RemoveImage(exampleImageName); err != nil {
		t.Errorf("Could not remove example image: %s", err.Error())
		return
	}
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

	// creating one and starting it up
	name, err := getPrefixedUUID(defaultTestContainerName)
	if err != nil {
		t.Errorf("Could not create prefixed container name: %s", err.Error())
		return
	}
	container, err := client.CreateContainer(name, defaultTestImage, []string{})
	if err != nil {
		t.Errorf("Could not create test container from image %s: %s", defaultTestImage, err.Error())
		return
	}
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

	// cleaning it up
	if err := client.RemoveContainer(container); err != nil {
		t.Errorf("Could not remove container %s: %s", container.Name, err.Error())
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
	name, err := getPrefixedUUID(defaultTestContainerName)
	if err != nil {
		t.Errorf("Could not create prefixed container name: %s", err.Error())
		return
	}
	container, err := client.CreateContainer(name, defaultDbImage, []string{})
	if err != nil {
		t.Errorf("Could not create test container from image %s: %s", defaultDbImage, err.Error())
		return
	}
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

	// cleaning it up
	if err := client.StopContainer(container); err != nil {
		t.Errorf("Could not stop container %s: %s", container.Name, err.Error())
		return
	}
	if err := client.RemoveContainer(container); err != nil {
		t.Errorf("Could not remove container %s: %s", container.Name, err.Error())
		return
	}
}
