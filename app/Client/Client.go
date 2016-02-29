package Client

import (
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
