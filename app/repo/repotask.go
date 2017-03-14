package repo

import (
	"github.com/ihsw/the-matrix/app/simpledocker"
	"github.com/ihsw/the-matrix/app/util"
)

type repoTask struct {
	repo Repo
	err  error
}

// NewRepos - creates a new list of repos from names
func NewRepos(names []string, simpleDocker simpledocker.SimpleDocker) ([]Repo, error) {
	// setting up the workers
	in := make(chan string)
	out := make(chan repoTask)
	worker := func() {
		for name := range in {
			repo, err := newRepo(name, simpleDocker)
			out <- repoTask{repo, err}
		}
	}
	postWork := func() { close(out) }
	util.Work(4, worker, postWork)

	// starting it up
	go func() {
		for _, name := range names {
			in <- name
		}
		close(in)
	}()

	// waiting for it to drain out
	repoTasks := []repoTask{}
	for repoTask := range out {
		repoTasks = append(repoTasks, repoTask)
	}

	// gathering up the list of repos and checking for errors
	repos := []Repo{}
	for _, repoTask := range repoTasks {
		if repoTask.err != nil {
			return []Repo{}, repoTask.err
		}

		repos = append(repos, repoTask.repo)
	}

	return repos, nil
}
