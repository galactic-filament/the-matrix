package endpoint

import (
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/galactic-filament/the-matrix/app/repo"
	"github.com/galactic-filament/the-matrix/app/resource"
	"github.com/galactic-filament/the-matrix/app/simpledocker"
	"github.com/galactic-filament/the-matrix/app/util"
)

type endpointTask struct {
	endpoint Endpoint
	err      error
}

// NewEndpoints - creates a new list of endpoints
func NewEndpoints(repos []repo.Repo, network *docker.Network, resources resource.Resources) ([]Endpoint, error) {
	// setting up the workers
	in := make(chan repo.Repo)
	out := make(chan endpointTask)
	worker := func() {
		for repo := range in {
			endpoint, err := NewEndpoint(repo, network, resources)
			out <- endpointTask{endpoint, err}
		}
	}
	postWork := func() { close(out) }
	util.Work(len(repos), worker, postWork)

	// starting it up
	go func() {
		for _, repo := range repos {
			in <- repo
		}
		close(in)
	}()

	// waiting for it to drain out
	endpointTasks := []endpointTask{}
	for endpointTask := range out {
		endpointTasks = append(endpointTasks, endpointTask)
	}

	// gathering up the list of endpoints and checking for errors
	endpoints := []Endpoint{}
	for _, endpointTask := range endpointTasks {
		if endpointTask.err != nil {
			return []Endpoint{}, endpointTask.err
		}

		endpoints = append(endpoints, endpointTask.endpoint)
	}

	return endpoints, nil
}

// CreateTestEndpointOpts - opts for corresponding func
type CreateTestEndpointOpts struct {
	Client   simpledocker.Client
	RepoName string
	Resource resource.Resource
	Network  *docker.Network
}

// CreateTestEndpoint - common test func for creating a test endpoint
func CreateTestEndpoint(opts CreateTestEndpointOpts) (Endpoint, error) {
	endpointRepo, err := repo.NewRepo(opts.RepoName, opts.Client)
	if err != nil {
		return Endpoint{}, err
	}

	e, err := NewEndpoint(
		endpointRepo,
		opts.Network,
		resource.Resources{Values: []resource.Resource{opts.Resource}},
	)
	if err != nil {
		return Endpoint{}, err
	}

	return e, nil
}

// CleanEndpoint - common test func used for cleaning up an endpoint
func CleanEndpoint(t *testing.T, e Endpoint) {
	if err := e.Clean(); err != nil {
		t.Errorf("Could not clean endpoint: %s", err.Error())
		return
	}
}
