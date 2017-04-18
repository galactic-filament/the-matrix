package repo

import (
	"fmt"

	"time"

	"github.com/ihsw/the-matrix/app/simpledocker"
)

const defaultRepoImageTag = "latest"

// GetImageName - returns formatted image ID of a repo
func GetImageName(name string) string { return fmt.Sprintf("ihsw/%s", name) }

// NewRepo - creates a new repo based on a local or remote docker image
func NewRepo(name string, client simpledocker.Client) (Repo, error) {
	r := Repo{name, client}
	imageName := GetImageName(r.Name)

	hasImage, err := client.HasImage(imageName)
	if err != nil {
		return Repo{}, err
	}

	if hasImage {
		return r, nil
	}

	err = client.PullImage(imageName, defaultRepoImageTag)
	if err != nil {
		return Repo{}, err
	}

	time.Sleep(5 * time.Second)

	return r, nil
}

// Repo - container to run tests with
type Repo struct {
	Name   string
	Client simpledocker.Client
}
