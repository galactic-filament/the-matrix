package client

import "github.com/ihsw/the-matrix/app/repo"

// NewClients - creates a list of new clients
func NewClients(repos []repo.Repo) ([]Client, error) {
	clients := []Client{}
	for _, repo := range repos {
		clients = append(clients, newClient(repo))
	}

	return clients, nil
}

func newClient(r repo.Repo) Client {
	return Client{r}
}

// Client - a container ran against an Endpoint
type Client struct {
	repo.Repo
}
