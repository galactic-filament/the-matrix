package main

import (
	"github.com/fsouza/go-dockerclient"
)

type repoTask struct {
	repo   repo
	err    error
	client *docker.Client
}

func runTasks(in chan repo, client *docker.Client) chan repoTask {
	out := make(chan repoTask)
	worker := func() {
		for repo := range in {
			out <- repoTask{repo: repo, err: repo.runTests()}
		}
	}
	postWork := func() { close(out) }
	work(2, worker, postWork)

	return out
}
