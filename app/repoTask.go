package main

type repoTask struct {
	repo repo
	err  error
}

func runTasks(in chan repo) chan repoTask {
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
