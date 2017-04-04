package simpledocker

import (
	"bytes"

	docker "github.com/fsouza/go-dockerclient"
)

// GetImage - fetches an image from id
func (c Client) GetImage(id string) (*docker.Image, error) {
	return c.dockerClient.InspectImage(id)
}

// HasImage - validates whether an image exists
func (c Client) HasImage(id string) (bool, error) {
	_, err := c.GetImage(id)
	if err != nil {
		if err != docker.ErrNoSuchImage {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// PullImage - pulls an image by repository and tag
func (c Client) PullImage(repository string, tag string) error {
	return c.dockerClient.PullImage(
		docker.PullImageOptions{Repository: repository, Tag: tag},
		docker.AuthConfiguration{},
	)
}

// RemoveImage - removes an image by id
func (c Client) RemoveImage(imageID string) error { return c.dockerClient.RemoveImage(imageID) }

// BuildImage - builds an image by name and dockerfile path
func (c Client) BuildImage(name string, contextDir string) error {
	outputBuffer := bytes.NewBuffer(nil)
	return c.dockerClient.BuildImage(docker.BuildImageOptions{
		Name:         name,
		ContextDir:   contextDir,
		OutputStream: outputBuffer,
	})
}
