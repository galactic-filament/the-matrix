package main

type repoTask struct {
	repo repo
	err  error
}

func runTasks(workerCount int, in chan repo) chan repoTask {
	out := make(chan repoTask)
	worker := func() {
		for repo := range in {
			out <- repoTask{repo: repo, err: repo.runTests()}
		}
	}
	postWork := func() { close(out) }
	work(workerCount, worker, postWork)

	return out
}
