package Work

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ihsw/the-matrix/app/Client"
	"github.com/ihsw/the-matrix/app/Endpoint"
)

// RunClient - starts up the Client container, runs it against the Endpoint, and exits
func RunClient(c Client.Client, e Endpoint.Endpoint) (*Client.TestOutput, error) {
	container, err := c.SimpleDocker.CreateContainer(
		fmt.Sprintf("%s-client", c.Name),
		fmt.Sprintf("ihsw/%s", c.Name),
	)
	if err != nil {
		return nil, err
	}

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

		return testOutput, errors.New("Test container exited with non-zero status")
	}

	err = c.SimpleDocker.RemoveContainer(container)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// RunEndpoint - starts up an Endpoint and runs Clients against it
func RunEndpoint(e Endpoint.Endpoint, clients []Client.Client) error {
	err := e.Start()
	if err != nil {
		return err
	}

	for _, c := range clients {
		testOutput, err := RunClient(c, e)
		if err != nil {
			log.WithFields(log.Fields{
				"endpoint": e.Repo.Name,
				"client":   c.Repo.Name,
			}).Warn("Client run failed")

			if testOutput != nil {
				for _, line := range testOutput.Results {
					log.WithFields(log.Fields{
						"endpoint": e.Repo.Name,
						"client":   c.Repo.Name,
						"expected": line.Expected,
						"actual":   line.Actual,
					}).Warn(line.Message)
				}
			}

			return err
		}
	}

	err = e.Stop()
	if err != nil {
		return err
	}

	err = e.RemoveContainers()
	if err != nil {
		return err
	}

	return nil
}
