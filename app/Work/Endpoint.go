package Work

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ihsw/the-matrix/app/Client"
	"github.com/ihsw/the-matrix/app/Endpoint"
	"github.com/ihsw/the-matrix/app/Resource"
	"github.com/ihsw/the-matrix/app/Util"
)

type endpointWorkTask struct {
	endpoint Endpoint.Endpoint
	err      error
}

// RunEndpoints - runs clients against endpoints
func RunEndpoints(endpoints []Endpoint.Endpoint, resources []Resource.Resource, clients []Client.Client) error {
	// setting up the workers
	in := make(chan Endpoint.Endpoint)
	out := make(chan endpointWorkTask)
	worker := func() {
		for endpoint := range in {
			err := runEndpoint(endpoint, resources, clients)
			out <- endpointWorkTask{
				endpoint: endpoint,
				err:      err,
			}
		}
	}
	postWork := func() { close(out) }
	Util.Work(len(endpoints), worker, postWork)

	// starting it up
	go func() {
		for _, endpoint := range endpoints {
			in <- endpoint
		}
		close(in)
	}()

	// waiting for it to drain out
	var lastError error
	for task := range out {
		if task.err == nil {
			continue
		}

		log.WithFields(log.Fields{
			"endpoint": task.endpoint.Repo.Name,
			"err":      task.err.Error(),
		}).Warn("Endpoint run failed")
		lastError = task.err
	}

	// cleaning the resources
	for _, resource := range resources {
		err := resource.Clean()
		if err == nil {
			continue
		}

		lastError = err
	}

	return lastError
}

func runEndpoint(e Endpoint.Endpoint, resources []Resource.Resource, clients []Client.Client) error {
	log.WithFields(log.Fields{
		"endpoint": e.Name,
	}).Info("Running endpoint")

	if err := runClients(e, clients); err != nil {
		return e.Clean(err)
	}

	return e.Clean(nil)
}
