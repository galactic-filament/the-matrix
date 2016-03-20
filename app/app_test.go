package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ihsw/the-matrix/app/Client"
	"github.com/ihsw/the-matrix/app/Endpoint"
	"github.com/ihsw/the-matrix/app/Repo"
	"github.com/ihsw/the-matrix/app/Resource"
	"github.com/ihsw/the-matrix/app/SimpleDocker"
	"github.com/ihsw/the-matrix/app/Work"
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"testing"
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
	simpleDocker, err := SimpleDocker.NewSimpleDocker("unix:///var/run/docker.sock")
	if err != nil {
		fail(t, err)
		return
	}

	// gathering up a list of resources
	resourceNames := map[string]string{
		"db": "Db",
	}
	resources, err := Resource.NewResources(simpleDocker, resourceNames)
	if err != nil {
		fail(t, err)
		return
	}

	// gathering up a list of endpoints
	endpointRepoNames := []string{
		"omega-jazz",
		// "pho-sho",
		// "go-home",
		// "py-lyfe",
		// "es-bueno",
		// "crazy-train",
		// "fur-elise",
	}
	endpointRepos, err := Repo.NewRepos(endpointRepoNames, simpleDocker)
	if err != nil {
		fail(t, err)
		return
	}
	endpoints, err := Endpoint.NewEndpoints(endpointRepos, resources)

	// gathering up a list of clients
	clientRepoNames := []string{
		"integration-nation",
	}
	clientRepos, err := Repo.NewRepos(clientRepoNames, simpleDocker)
	if err != nil {
		fail(t, err)
		return
	}
	clients, err := Client.NewClients(clientRepos)
	if err != nil {
		fail(t, err)
	}

	err = Work.RunEndpoints(endpoints, resources, clients)
	if err != nil {
		fail(t, err)
		return
	}
}

func fail(t *testing.T, err error) {
	assert.Equal(t, err.Error(), nil)
}
