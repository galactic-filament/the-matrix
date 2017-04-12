package simpledocker

import (
	"testing"

	docker "github.com/fsouza/go-dockerclient"
)

// CleanupContainer - used for deferred cleanup of containers
func CleanupContainer(t *testing.T, client Client, container *docker.Container) {
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

// CleanupImage - used for deferred cleanup of images
func CleanupImage(t *testing.T, client Client, imageID string) {
	if err := client.RemoveImage(imageID); err != nil {
		t.Errorf("Could not remove image: %s", err.Error())
		return
	}
}

// CleanupNetwork - used for deferred cleanup of networks
func CleanupNetwork(t *testing.T, client Client, network *docker.Network) {
	if err := client.RemoveNetwork(network); err != nil {
		t.Errorf("Could not remove network: %s", err.Error())
		return
	}
}
