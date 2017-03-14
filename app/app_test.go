package main

import (
	"os"
	"runtime"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/ihsw/the-matrix/app/client"
	"github.com/ihsw/the-matrix/app/endpoint"
	"github.com/ihsw/the-matrix/app/repo"
	"github.com/ihsw/the-matrix/app/resource"
	"github.com/ihsw/the-matrix/app/simpledocker"
	"github.com/ihsw/the-matrix/app/work"
	"github.com/stretchr/testify/assert"
)

func init() {
	logLevel := log.InfoLevel
	if os.Getenv("ENV") == "travis" {
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func TestTestSuite(t *testing.T) {
	// connecting to docker
	simpleDocker, err := simpledocker.NewSimpleDocker("unix:///var/run/docker.sock")
	if err != nil {
		fail(t, err)
		return
	}

	// gathering up a list of resources
	resourceNames := map[string]string{
		"db": "Db",
	}
	resources, err := resource.NewResources(simpleDocker, resourceNames)
	if err != nil {
		fail(t, err)
		return
	}

	// gathering up a list of endpoints
	endpointRepoNames := []string{
		// "omega-jazz",
		// "pho-sho",
		// "go-home",
		"py-lyfe",
		// "es-bueno",
		// "crazy-train",
		// "fur-elise",
	}
	endpointRepos, err := repo.NewRepos(endpointRepoNames, simpleDocker)
	if err != nil {
		fail(t, err)
		return
	}
	endpoints, err := endpoint.NewEndpoints(endpointRepos, resources)

	// gathering up a list of clients
	clientRepoNames := []string{
		"integration-nation",
	}
	clientRepos, err := repo.NewRepos(clientRepoNames, simpleDocker)
	if err != nil {
		fail(t, err)
		return
	}
	clients, err := client.NewClients(clientRepos)
	if err != nil {
		fail(t, err)
	}

	err = work.RunEndpoints(endpoints, resources, clients)
	if err != nil {
		fail(t, err)
		return
	}
}

func fail(t *testing.T, err error) {
	assert.Equal(t, err.Error(), nil)
}
