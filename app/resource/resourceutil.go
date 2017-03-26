package resource

import "github.com/ihsw/the-matrix/app/simpledocker"

// NewResources - generates a new list of resources
func NewResources(client simpledocker.Client, optList []Opts) ([]Resource, error) {
	resources := []Resource{}
	for _, opts := range optList {
		resource, err := NewResource(client, opts)
		if err != nil {
			return []Resource{}, err
		}

		resources = append(resources, resource)
	}

	return resources, nil
}
