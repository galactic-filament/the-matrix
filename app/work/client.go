package work

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/client"
	"github.com/ihsw/the-matrix/app/endpoint"
	"github.com/ihsw/the-matrix/app/util"
)

type clientWorkTask struct {
	endpoint   endpoint.Endpoint
	client     client.Client
	testOutput *client.TestOutput
	err        error
}

func runClients(endpoint endpoint.Endpoint, clients []client.Client) error {
	// setting up the workers
	in := make(chan client.Client)
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
	util.Work(len(clients), worker, postWork)

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

		log.WithFields(log.Fields{
			"endpoint": task.endpoint.Repo.Name,
			"client":   task.client.Repo.Name,
			"err":      task.err.Error(),
		}).Warn("Client run failed")
		lastError = task.err

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

func runClient(c client.Client, e endpoint.Endpoint) (*client.TestOutput, error) {
	log.WithFields(log.Fields{
		"endpoint": e.Name,
		"client":   c.Name,
	}).Info("Running client")

	endpointHostname := "ApiServer"
	clientContainer, err := c.SimpleDocker.CreateContainer(
		fmt.Sprintf("%s-%s-client", e.Name, c.Name),
		fmt.Sprintf("ihsw/%s", c.Name),
		[]string{fmt.Sprintf("API_HOST=%s", endpointHostname)},
	)
	if err != nil {
		return nil, err
	}

	failed, err := c.SimpleDocker.RunContainer(clientContainer, []string{
		fmt.Sprintf("%s:%s", e.Container.ID, endpointHostname),
	})
	if err != nil {
		return nil, cleanClient(c, clientContainer, err)
	}

	containerLogs, err := c.SimpleDocker.GetContainerLogs(clientContainer)
	if err != nil {
		return nil, err
	}

	if failed {
		testOutput, err := client.ParseClientLogs(containerLogs)
		if err != nil {
			return nil, errors.New("Client logs could not be parsed")
		}

		return testOutput, cleanClient(c, clientContainer, errors.New("Test container exited with non-zero status"))
	}

	return nil, cleanClient(c, clientContainer, nil)
}

func cleanClient(c client.Client, container *docker.Container, prevErr error) error {
	return prevErr
	// err := c.SimpleDocker.RemoveContainer(container)
	// if err != nil {
	// 	return err
	// }

	// return prevErr
}
