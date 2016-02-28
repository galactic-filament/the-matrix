package Work

import (
	"errors"
	"fmt"
	"github.com/ihsw/the-matrix/app/Client"
	"github.com/ihsw/the-matrix/app/Endpoint"
)

// RunClient - starts up the Client container, runs it against the Endpoint, and exits
func RunClient(c Client.Client, e Endpoint.Endpoint) (string, error) {
	container, err := c.SimpleDocker.CreateContainer(
		fmt.Sprintf("%s-client", c.Name),
		fmt.Sprintf("ihsw/%s", c.Name),
	)
	if err != nil {
		return "", err
	}

	endpointContainerID, err := e.GetContainerID()
	if err != nil {
		return "", err
	}

	failed, err := c.SimpleDocker.RunContainer(container, endpointContainerID)
	if err != nil {
		return "", err
	}

	containerLogs, err := c.SimpleDocker.GetContainerLogs(container)
	if err != nil {
		return "", err
	}

	if failed {
		return containerLogs, errors.New("Test container exited with non-zero status")
	}

	err = c.SimpleDocker.RemoveContainer(container)
	if err != nil {
		return containerLogs, err
	}

	return containerLogs, nil
}

// RunEndpoint - starts up an Endpoint and runs Clients against it
func RunEndpoint(e Endpoint.Endpoint, clients []Client.Client) error {
	err := e.Start()
	if err != nil {
		return err
	}

	for _, c := range clients {
		_, err := RunClient(c, e)
		if err != nil {
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
