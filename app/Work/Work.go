package Work

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/Client"
	"github.com/ihsw/the-matrix/app/Endpoint"
	"time"
)

// runClient - starts up the Client container, runs it against the Endpoint, and exits
func runClient(c Client.Client, e Endpoint.Endpoint) (*Client.TestOutput, error) {
	log.WithFields(log.Fields{
		"endpoint": e.Name,
		"client":   c.Name,
	}).Info("Running client")

	container, err := c.SimpleDocker.CreateContainer(
		fmt.Sprintf("%s-%s-client", e.Name, c.Name),
		fmt.Sprintf("ihsw/%s", c.Name),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint": e.Name,
			"client":   c.Name,
			"err":      err.Error(),
		}).Warn("Could not create a client container")

		return nil, err
	}

	time.Sleep(10 * time.Second)

	endpointContainerID, err := e.GetContainerID()
	if err != nil {
		return nil, err
	}

	failed, err := c.SimpleDocker.RunContainer(container, endpointContainerID)
	if err != nil {
		return nil, err
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

		return testOutput, cleanClient(container, c, errors.New("Test container exited with non-zero status"))
	}

	return nil, cleanClient(container, c, nil)
}

// cleanClient - removes the client container
func cleanClient(container *docker.Container, c Client.Client, prevErr error) error {
	err := c.SimpleDocker.RemoveContainer(container)
	if err != nil {
		return err
	}

	return prevErr
}

// runEndpoint - starts up an Endpoint and runs Clients against it
func runEndpoint(e Endpoint.Endpoint, clients []Client.Client) error {
	log.WithFields(log.Fields{
		"endpoint": e.Name,
	}).Info("Running endpoint")

	err := e.Start()
	if err != nil {
		return err
	}

	return cleanEndpoint(e, runClients(e, clients))
}

// cleanEndpoint - stops all endpoint services and removes their containers
func cleanEndpoint(e Endpoint.Endpoint, prevErr error) error {
	err := e.Stop()
	if err != nil {
		return err
	}

	err = e.RemoveContainers()
	if err != nil {
		return err
	}

	return prevErr
}
