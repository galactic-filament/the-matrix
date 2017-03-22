package resource

import "github.com/ihsw/the-matrix/app/simpledocker"

// NewResources - generates a new list of resources
func NewResources(client simpledocker.Client, names []string) ([]Resource, error) {
	resources := []Resource{}
	for _, name := range names {
		resource, err := newResource(name, client)
		if err != nil {
			return []Resource{}, err
		}

		resources = append(resources, resource)
	}

	return resources, nil
}
