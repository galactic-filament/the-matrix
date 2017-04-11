package resource

import "github.com/ihsw/the-matrix/app/simpledocker"

// NewResources - generates a new list of resources
func NewResources(client simpledocker.Client, optList []Opts) (Resources, error) {
	resources := Resources{[]Resource{}}
	for _, opts := range optList {
		resource, err := NewResource(client, opts)
		if err != nil {
			return Resources{}, err
		}

		resources.Values = append(resources.Values, resource)
	}

	return resources, nil
}

// Resources - a list of resources
type Resources struct {
	Values []Resource
}

// Clean - cleans all resources
func (r Resources) Clean() error {
	for _, resource := range r.Values {
		if err := resource.Clean(); err != nil {
			return err
		}
	}

	return nil
}
