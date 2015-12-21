package main

import (
	"fmt"
	"os"
)

// repo manager
type repoManager struct {
	repos            []repo
	cloneDestination string
}

func (r repoManager) cleanup() error {
	return os.RemoveAll(r.cloneDestination)
}

// repo
type repo struct {
	manager repoManager
	name    string
}

func (r repo) clonePath() string {
	return fmt.Sprintf("%s/%s", r.manager.cloneDestination, r.name)
}

func (r repo) cleanup() error {
	return os.RemoveAll(r.clonePath())
}

func (r repo) runTests() error {
	cloneCommand := fmt.Sprintf(
		"git clone %s %s",
		fmt.Sprintf("https://github.com/ihsw/%s.git", r.name),
		r.clonePath(),
	)
	if err := runCommand(cloneCommand); err != nil {
		if err := r.cleanup(); err != nil {
			return err
		}
		return err
	}

	return r.cleanup()
}
