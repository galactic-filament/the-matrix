package client

import (
	"github.com/ihsw/the-matrix/app/Repo"
)

// NewClients - creates a list of new clients
func NewClients(repos []Repo.Repo) ([]Client, error) {
	clients := []Client{}
	for _, repo := range repos {
		clients = append(clients, newClient(repo))
	}

	return clients, nil
}

func newClient(r Repo.Repo) Client {
	return Client{r}
}

// Client - a container ran against an Endpoint
type Client struct {
	Repo.Repo
}
