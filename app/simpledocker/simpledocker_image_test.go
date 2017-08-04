package simpledocker

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/util"
)

const defaultTestImageOutput = "hello world\n"
const defaultTestImage = "hello-world"
const defaultTestImageTag = "latest"
const defaultDbImage = "postgres"

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

	fmt.Printf("Pulling image %s\n", defaultTestImage)
	if err := client.PullImage(defaultTestImage, defaultTestImageTag); err != nil {
		t.Errorf("Could not pull test image %s: %s", defaultTestImage, err.Error())
	}
	CleanupImage(t, client, defaultTestImage)
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
		fmt.Printf("Pulling image %s\n", defaultTestImage)
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
	CleanupImage(t, client, exampleImageName)
}
