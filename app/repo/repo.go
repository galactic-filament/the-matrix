package repo

import (
	"fmt"

	"github.com/ihsw/the-matrix/app/simpledocker"
)

const defaultRepoImageTag = "latest"

func newRepo(name string, client simpledocker.Client) (Repo, error) {
	r := Repo{name, client}

	hasImage, err := client.HasImage(r.getImageID())
	if err != nil {
		return Repo{}, err
	}

	if hasImage {
		return r, nil
	}

	err = client.PullImage(r.getImageID(), defaultRepoImageTag)
	if err != nil {
		return Repo{}, err
	}

	return r, nil
}

// Repo - container to run tests with
type Repo struct {
	Name   string
	Client simpledocker.Client
}

func (r Repo) getImageID() string { return fmt.Sprintf("ihsw/%s", r.Name) }
