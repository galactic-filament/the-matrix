package endpoint

import (
	"github.com/ihsw/the-matrix/app/Repo"
	"github.com/ihsw/the-matrix/app/Resource"
)

type endpointTask struct {
	endpoint Endpoint
	err      error
}

// NewEndpoints - creates a new list of endpoints
func NewEndpoints(repos []Repo.Repo, resources []Resource.Resource) ([]Endpoint, error) {
	// setting up the workers
	in := make(chan Repo.Repo)
	out := make(chan endpointTask)
	worker := func() {
		for repo := range in {
			endpoint, err := newEndpoint(repo, resources)
			out <- endpointTask{endpoint, err}
		}
	}
	postWork := func() { close(out) }
	Util.Work(len(repos), worker, postWork)

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
