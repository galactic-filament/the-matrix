package repo

import (
	"fmt"

	"github.com/ihsw/the-matrix/app/simpledocker"
)

const defaultRepoImageTag = "latest"

// GetImageID - returns formatted image ID of a repo
func GetImageID(name string) string { return fmt.Sprintf("ihsw/%s", name) }

// NewRepo - creates a new repo based on a local or remote docker image
func NewRepo(name string, client simpledocker.Client) (Repo, error) {
	r := Repo{name, client}
	imageID := GetImageID(r.Name)

	hasImage, err := client.HasImage(imageID)
	if err != nil {
		return Repo{}, err
	}

	if hasImage {
		return r, nil
	}

	err = client.PullImage(imageID, defaultRepoImageTag)
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
