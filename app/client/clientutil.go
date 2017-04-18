package client

import "testing"

// DefaultTestClientName - common const for testing clients
const DefaultTestClientName = "integration-nation"

// CleanClient - common func used for cleaning up a client
func CleanClient(t *testing.T, c Client) {
	if err := c.Clean(); err != nil {
		t.Errorf("Could not clean client: %s", err.Error())
		return
	}
}
