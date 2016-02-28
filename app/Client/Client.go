package Client

import (
	"errors"
	"fmt"
	"github.com/ihsw/the-matrix/app/Endpoint"
	"github.com/ihsw/the-matrix/app/Repo"
	"github.com/ihsw/the-matrix/app/SimpleDocker"
)

// Client - a container ran against an Endpoint
type Client struct {
	Repo.Repo
}

// NewClient - creates a new Client
func NewClient(r Repo.Repo) Client {
	return Client{r}
}

// NewClients - creates a list of new clients
func NewClients(repoNames []string, gitFormat string, cloneDirectory string, simpleDocker SimpleDocker.SimpleDocker) ([]Client, error) {
	// generating a list of new repos
	repos, err := Repo.NewRepos(repoNames, gitFormat, cloneDirectory, simpleDocker)
	if err != nil {
		return []Client{}, err
	}

	// generating a list of new clients
	clients := []Client{}
	for _, repo := range repos {
		clients = append(clients, NewClient(repo))
	}

	return clients, nil
}

// Run - starts up the Client container, runs it against the Endpoint, and exits
func (c Client) Run(e Endpoint.Endpoint) (string, error) {
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
