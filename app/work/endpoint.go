package work

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ihsw/the-matrix/app/client"
	"github.com/ihsw/the-matrix/app/endpoint"
	"github.com/ihsw/the-matrix/app/resource"
	"github.com/ihsw/the-matrix/app/util"
)

type endpointWorkTask struct {
	endpoint endpoint.Endpoint
	err      error
}

// RunEndpoints - runs clients against endpoints
func RunEndpoints(endpoints []endpoint.Endpoint, resources []resource.Resource, clients []client.Client) error {
	// setting up the workers
	in := make(chan endpoint.Endpoint)
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
	util.Work(len(endpoints), worker, postWork)

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
	// for _, resource := range resources {
	// 	err := resource.Clean()
	// 	if err == nil {
	// 		continue
	// 	}

	// 	lastError = err
	// }

	return lastError
}

func runEndpoint(e endpoint.Endpoint, resources []resource.Resource, clients []client.Client) error {
	log.WithFields(log.Fields{
		"endpoint": e.Name,
	}).Info("Running endpoint")

	if err := runClients(e, clients); err != nil {
		return err
		// return e.Clean(err)
	}

	return nil
	// return e.Clean(nil)
}
