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

// GetLinkLineList - returns a list of docker link lines
func (r Resources) GetLinkLineList() []string {
	linkLineList := []string{}
	for _, resource := range r.Values {
		linkLineList = append(linkLineList, resource.GetLinkLine())
	}

	return linkLineList
}

// GetEnvVarsList - returns a list of env vars for a group of resources
func (r Resources) GetEnvVarsList() []string {
	envVarsList := []string{}
	for _, resource := range r.Values {
		for _, envVars := range resource.GetEnvVars() {
			envVarsList = append(envVarsList, envVars)
		}
	}

	return envVarsList
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
