package simpledocker

import (
	docker "github.com/fsouza/go-dockerclient"
)

// Client - simplified interface to the docker client
type Client struct {
	dockerClient *docker.Client
}

// NewClient - creates a new SimpleDocker
func NewClient(dockerClient *docker.Client) Client {
	return Client{dockerClient: dockerClient}
}
