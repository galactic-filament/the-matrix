package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ihsw/the-matrix/app/Client"
	"github.com/ihsw/the-matrix/app/Endpoint"
	"github.com/ihsw/the-matrix/app/SimpleDocker"
	"github.com/ihsw/the-matrix/app/Work"
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"testing"
)

func fail(t *testing.T, err error) {
	assert.Equal(t, err.Error(), nil)
}

func init() {
	logLevel := log.WarnLevel
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

	// misc
	gitFormat := "https://github.com/ihsw/%s.git"

	// gathering up a list of clients
	clientRepoNames := []string{
		"integration-nation",
	}
	clients, err := Client.NewClients(
		clientRepoNames,
		gitFormat,
		"./client-repos",
		simpleDocker,
	)
	if err != nil {
		fail(t, err)
		return
	}

	// gathering up a list of endpoints
	endpointRepoNames := []string{
		"omega-jazz",
		"pho-sho",
		"go-home",
		"py-lyfe",
		"es-bueno",
		"crazy-train",
		"fur-elise",
	}
	endpoints, err := Endpoint.NewEndpoints(
		endpointRepoNames,
		gitFormat,
		"./endpoint-repos",
		simpleDocker,
	)
	if err != nil {
		fail(t, err)
		return
	}

	err = Work.RunEndpoints(endpoints, clients)
	if err != nil {
		fail(t, err)
		return
	}
}
