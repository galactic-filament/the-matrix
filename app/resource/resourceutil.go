package resource

import "testing"

// CleanResource - common test func used for cleaning up a resource
func CleanResource(t *testing.T, resource Resource) {
	if err := resource.Clean(); err != nil {
		t.Errorf("Could not clean resource %s: %s", resource.name, err.Error())
		return
	}
}
