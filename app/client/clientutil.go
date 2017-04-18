package client

import (
	"testing"

	docker "github.com/fsouza/go-dockerclient"
)

// DefaultTestClientName - common const for testing clients
const DefaultTestClientName = "integration-nation"

// CleanClient - common func used for cleaning up a client
func CleanClient(t *testing.T, c Client, container *docker.Container) {
	if err := c.Clean(container); err != nil {
		t.Errorf("Could not clean client: %s", err.Error())
		return
	}
}
