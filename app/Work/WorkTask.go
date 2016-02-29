package Work

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ihsw/the-matrix/app/Client"
	"github.com/ihsw/the-matrix/app/Endpoint"
	"github.com/ihsw/the-matrix/app/Util"
)

type workTask struct {
	endpoint   Endpoint.Endpoint
	client     Client.Client
	testOutput *Client.TestOutput
	err        error
}

func runClients(endpoint Endpoint.Endpoint, clients []Client.Client) error {
	// setting up the workers
	in := make(chan Client.Client)
	out := make(chan workTask)
	worker := func() {
		for client := range in {
			testOutput, err := runClient(client, endpoint)
			out <- workTask{
				endpoint:   endpoint,
				client:     client,
				err:        err,
				testOutput: testOutput,
			}
		}
	}
	postWork := func() { close(out) }

	// starting it up
	go func() {
		for _, client := range clients {
			in <- client
		}
		close(in)
	}()
	Util.Work(4, worker, postWork)

	// waiting for it to drain out
	var lastError error
	for workTask := range out {
		if workTask.err != nil {
			continue
		}

		lastError = workTask.err
		log.WithFields(log.Fields{
			"endpoint": workTask.endpoint.Repo.Name,
			"client":   workTask.client.Repo.Name,
			"err":      workTask.err.Error(),
		}).Warn("Client run failed")

		if workTask.testOutput != nil {
			for _, line := range workTask.testOutput.Results {
				log.WithFields(log.Fields{
					"endpoint": workTask.endpoint.Repo.Name,
					"client":   workTask.client.Repo.Name,
					"expected": line.Expected,
					"actual":   line.Actual,
				}).Warn(line.Message)
			}
		}
	}

	return lastError
}
