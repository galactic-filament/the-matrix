package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

func fail(t *testing.T, err error) {
	assert.Equal(t, err.Error(), nil)
}

func init() {
	log.SetLevel(log.WarnLevel)
}

func TestTestSuite(t *testing.T) {
	// misc
	repoManager := repoManager{cloneDestination: "./repos"}
	repoNames := []string{
		"omega-jazz",
		"pho-sho",
		"go-home",
		"py-lyfe",
		"es-bueno",
	}
	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		fail(t, err)
		return
	}

	// spinning up some workers
	reposToDo := make(chan repo)
	reposFinished := runTasks(reposToDo, client)

	// starting it up
	go func() {
		for _, repoName := range repoNames {
			reposToDo <- repo{manager: repoManager, name: repoName, client: client}
		}
		close(reposToDo)
	}()

	// waiting for it to drain out
	tasksFinished := make([]repoTask, len(repoNames))
	for task := range reposFinished {
		tasksFinished = append(tasksFinished, task)
	}

	// going over the list of tasks and dumping out debug info where appropriate
	for _, task := range tasksFinished {
		if err := task.err; err != nil {
			fail(t, err)
		}
	}

	if err := repoManager.cleanup(); err != nil {
		fail(t, err)
		return
	}
}
