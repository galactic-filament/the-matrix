package Endpoint

import (
	"github.com/ihsw/the-matrix/app/Repo"

	"github.com/ihsw/the-matrix/app/SimpleDocker"
	"strings"
)

// Endpoint - a container ran against an Endpoint
type Endpoint struct {
	Repo.Repo
}

// NewEndpoint - creates a new Endpoint
func NewEndpoint(r Repo.Repo) Endpoint {
	return Endpoint{r}
}

// NewEndpoints - creates a list of new endpoints
func NewEndpoints(repoNames []string, gitFormat string, cloneDirectory string, simpleDocker SimpleDocker.SimpleDocker) ([]Endpoint, error) {
	// generating a list of new repos
	repos, err := Repo.NewRepos(repoNames, gitFormat, cloneDirectory, simpleDocker)
	if err != nil {
		return []Endpoint{}, err
	}

	// generating a list of new endpoints
	endpoints := []Endpoint{}
	for _, repo := range repos {
		endpoints = append(endpoints, NewEndpoint(repo))
	}

	return endpoints, nil
}

// Start - starts up an endpoint
func (e Endpoint) Start() error {
	_, err := e.RunCommand("docker-compose up -d web-test")
	if err != nil {
		return err
	}

	return nil
}

// Stop - stops an endpoint
func (e Endpoint) Stop() error {
	_, err := e.RunCommand("docker-compose stop")
	if err != nil {
		return err
	}

	return nil
}

// RemoveContainers - removes containers associated with an endpoint
func (e Endpoint) RemoveContainers() error {
	_, err := e.RunCommand("docker-compose rm -v --force")
	if err != nil {
		return err
	}

	return nil
}

// GetContainerID - gets this endpoint's container ID
func (e Endpoint) GetContainerID() (string, error) {
	output, err := e.RunCommand("docker-compose ps -q web-test")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
