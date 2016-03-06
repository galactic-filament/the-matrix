package Repo

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ihsw/the-matrix/app/SimpleDocker"
	"github.com/ihsw/the-matrix/app/Util"
	"os"
)

// Repo - container to run tests with
type Repo struct {
	Name           string
	gitFormat      string
	cloneDirectory string
	SimpleDocker   SimpleDocker.SimpleDocker
}

func newRepo(name string, gitFormat string, cloneDirectory string, simpleDocker SimpleDocker.SimpleDocker) (Repo, error) {
	log.WithFields(log.Fields{
		"name": name,
	}).Info("Creating new repo")

	r := Repo{
		Name:           name,
		gitFormat:      gitFormat,
		cloneDirectory: cloneDirectory,
		SimpleDocker:   simpleDocker,
	}

	err := r.clone()
	if err != nil {
		return Repo{}, err
	}

	err = r.buildImages()
	if err != nil {
		return Repo{}, err
	}

	return r, nil
}

// Clone - clones the repo
func (r Repo) clone() error {
	cmd := fmt.Sprintf("git clone %s %s", r.gitURL(), r.clonePath())
	if _, err := Util.RunCommand(cmd); err != nil {
		return err
	}

	return nil
}

func (r Repo) clonePath() string {
	return fmt.Sprintf("%s/%s", r.cloneDirectory, r.Name)
}

func (r Repo) gitURL() string {
	return fmt.Sprintf(r.gitFormat, r.Name)
}

// RunCommand - runs a shell command in the context of this repo
func (r Repo) RunCommand(cmd string) ([]byte, error) {
	return Util.RunCommand(fmt.Sprintf("cd %s && %s", r.clonePath(), cmd))
}

// BuildImages - runs the build-images command found in all repos
func (r Repo) buildImages() error {
	cmd := "./bin/build-images"
	if _, err := r.RunCommand(cmd); err != nil {
		return err
	}

	return nil
}

// RemoveDir - removes the contents of the clone path
func (r Repo) RemoveDir() error {
	err := os.RemoveAll(r.clonePath())
	if err != nil {
		return err
	}

	return nil
}
