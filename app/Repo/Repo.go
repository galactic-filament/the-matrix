package Repo

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ihsw/the-matrix/app/SimpleDocker"
	"github.com/ihsw/the-matrix/app/Util"
)

// Repo - container to run tests with
type Repo struct {
	Name         string
	SimpleDocker SimpleDocker.SimpleDocker
}

func newRepo(name string, simpleDocker SimpleDocker.SimpleDocker) (Repo, error) {
	log.WithFields(log.Fields{
		"name": name,
	}).Info("Creating new repo")

	r := Repo{
		Name:         name,
		SimpleDocker: simpleDocker,
	}

	repoName := fmt.Sprintf("ihsw/%s", name)
	if err := r.pullImage(repoName); err != nil {
		log.WithFields(log.Fields{
			"name":     name,
			"err":      err.Error(),
			"repoName": repoName,
		}).Warn("Could not pull image")

		return Repo{}, err
	}

	log.WithFields(log.Fields{
		"name": name,
	}).Info("Repo create success")

	return r, nil
}

func (r Repo) pullImage(repoName string) error {
	cmd := fmt.Sprintf("docker pull %s", repoName)
	if _, err := Util.RunCommand(cmd); err != nil {
		return err
	}

	return nil
}
