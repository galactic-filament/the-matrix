package Work

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
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

func runClient(c Client.Client, e Endpoint.Endpoint) (*Client.TestOutput, error) {
	log.WithFields(log.Fields{
		"endpoint": e.Name,
		"client":   c.Name,
	}).Info("Running client")

	container, err := c.SimpleDocker.CreateContainer(
		fmt.Sprintf("%s-%s-client", e.Name, c.Name),
		fmt.Sprintf("ihsw/%s", c.Name),
		[]string{},
	)
	if err != nil {
		return nil, err
	}

	failed, err := c.SimpleDocker.RunContainer(container, []string{fmt.Sprintf("%s:ApiServer", e.Container.ID)})
	if err != nil {
		return nil, cleanClient(c, container, err)
	}

	containerLogs, err := c.SimpleDocker.GetContainerLogs(container)
	if err != nil {
		return nil, err
	}

	if failed {
		testOutput, err := Client.ParseClientLogs(containerLogs)
		if err != nil {
			return nil, errors.New("Client logs could not be parsed")
		}

		return testOutput, cleanClient(c, container, errors.New("Test container exited with non-zero status"))
	}

	return nil, cleanClient(c, container, nil)
}

func cleanClient(c Client.Client, container *docker.Container, prevErr error) error {
	err := c.SimpleDocker.RemoveContainer(container)
	if err != nil {
		return err
	}

	return prevErr
}
