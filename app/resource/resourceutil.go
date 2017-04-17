package resource

import (
	"testing"
)

// DefaultTestResourceName - common resource name for testing
const DefaultTestResourceName = "db"

// CleanResource - common test func used for cleaning up a resource
func CleanResource(t *testing.T, resource Resource) {
	if err := resource.Clean(); err != nil {
		t.Errorf("Could not clean resource %s: %s", resource.name, err.Error())
		return
	}
}

// CleanResources - common test func used for cleaning up resources
func CleanResources(t *testing.T, resources Resources) {
	if err := resources.Clean(); err != nil {
		t.Errorf("Could not clean resources: %s", err.Error())
		return
	}
}
