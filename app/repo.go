package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
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

func (r repo) runCommand(name string) error {
	if err := runCommand(name); err != nil {
		if err := r.cleanup(); err != nil {
			return err
		}
		return err
	}
	return nil
}

func (r repo) runRepoCommand(name string) error {
	return r.runCommand(fmt.Sprintf("cd %s && %s", r.clonePath(), name))
}

func (r repo) runTests() error {
	// cloning the repo in
	log.WithFields(log.Fields{"repo": r.name}).Info("Cloning")
	cloneCommand := fmt.Sprintf(
		"git clone %s %s",
		fmt.Sprintf("https://github.com/ihsw/%s.git", r.name),
		r.clonePath(),
	)
	if err := r.runCommand(cloneCommand); err != nil {
		return err
	}

	// building up the related images
	log.WithFields(log.Fields{"repo": r.name}).Info("Building images")
	if err := r.runRepoCommand("./bin/build-images"); err != nil {
		return err
	}

	// starting up the web-test service
	log.WithFields(log.Fields{"repo": r.name}).Info("Up web-test")
	if err := r.runRepoCommand("docker-compose up -d web-test"); err != nil {
		return err
	}

	// TODO: running the tests against the web-test service container

	return r.cleanup()
}

func (r repo) cleanup() error {
	// TODO: clean up the test container

	// stopping the services
	log.WithFields(log.Fields{"repo": r.name}).Info("Stop all")
	if err := r.runRepoCommand("docker-compose stop"); err != nil {
		return err
	}

	// cleaning up the web-test service containers
	log.WithFields(log.Fields{"repo": r.name}).Info("Remove containers")
	if err := r.runRepoCommand("docker rm -v $(docker-compose ps -q)"); err != nil {
		return err
	}

	// removing the cloned repo
	log.WithFields(log.Fields{"repo": r.name}).Info("Remove repo files")
	if err := os.RemoveAll(r.clonePath()); err != nil {
		return err
	}

	return nil
}
