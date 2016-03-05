package Work

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ihsw/the-matrix/app/Client"
	"github.com/ihsw/the-matrix/app/Endpoint"
	"github.com/ihsw/the-matrix/app/Util"
)

type clientWorkTask struct {
	endpoint   Endpoint.Endpoint
	client     Client.Client
	testOutput *Client.TestOutput
	err        error
}

func runClients(endpoint Endpoint.Endpoint, clients []Client.Client) error {
	// setting up the workers
	in := make(chan Client.Client)
	out := make(chan clientWorkTask)
	worker := func() {
		for client := range in {
			testOutput, err := runClient(client, endpoint)
			out <- clientWorkTask{
				endpoint:   endpoint,
				client:     client,
				err:        err,
				testOutput: testOutput,
			}
		}
	}
	postWork := func() { close(out) }
	Util.Work(len(clients), worker, postWork)

	// starting it up
	go func() {
		for _, client := range clients {
			in <- client
		}
		close(in)
	}()

	// waiting for it to drain out
	var lastError error
	for task := range out {
		if task.err == nil {
			continue
		}

		lastError = task.err
		log.WithFields(log.Fields{
			"endpoint": task.endpoint.Repo.Name,
			"client":   task.client.Repo.Name,
			"err":      task.err.Error(),
		}).Warn("Client run failed")

		if task.testOutput != nil {
			for _, line := range task.testOutput.Results {
				log.WithFields(log.Fields{
					"endpoint": task.endpoint.Repo.Name,
					"client":   task.client.Repo.Name,
					"expected": line.Expected,
					"actual":   line.Actual,
				}).Warn(line.Message)
			}
		}
	}

	return lastError
}

type endpointWorkTask struct {
	endpoint Endpoint.Endpoint
	err      error
}

// RunEndpoints - runs clients against endpoints
func RunEndpoints(endpoints []Endpoint.Endpoint, clients []Client.Client) error {
	// setting up the workers
	in := make(chan Endpoint.Endpoint)
	out := make(chan endpointWorkTask)
	worker := func() {
		for endpoint := range in {
			err := runEndpoint(endpoint, clients)
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
		}).Warn("Endpoint run failed")
		lastError = task.err
	}

	return lastError
}
