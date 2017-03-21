package repo

import (
	"fmt"

	"github.com/ihsw/the-matrix/app/simpledocker"
)

const defaultRepoImageTag = "latest"

// GetImageID - returns formatted image ID of a repo
func GetImageID(name string) string { return fmt.Sprintf("ihsw/%s", name) }

func newRepo(name string, client simpledocker.Client) (Repo, error) {
	r := Repo{name, client}

	hasImage, err := client.HasImage(GetImageID(r.Name))
	if err != nil {
		return Repo{}, err
	}

	if hasImage {
		return r, nil
	}

	err = client.PullImage(GetImageID(r.Name), defaultRepoImageTag)
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
