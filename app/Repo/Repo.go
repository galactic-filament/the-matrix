package Repo

import (
	"fmt"
	"time"

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
	imageID := fmt.Sprintf("ihsw/%s", name)
	r := Repo{
		Name:         name,
		SimpleDocker: simpleDocker,
	}
	if _, err := simpleDocker.GetImage(imageID); err == nil {
		return r, nil
	}

	log.WithFields(log.Fields{
		"name": name,
	}).Info("Creating new repo")
	startTime := time.Now()
	if err := r.pullImage(imageID); err != nil {
		log.WithFields(log.Fields{
			"name":     name,
			"err":      err.Error(),
			"repoName": imageID,
		}).Warn("Could not pull image")

		return Repo{}, err
	}

	log.WithFields(log.Fields{
		"name":     name,
		"duration": fmt.Sprintf("%v", time.Now().Sub(startTime)),
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
