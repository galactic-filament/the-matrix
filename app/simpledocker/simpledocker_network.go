package simpledocker

import docker "github.com/fsouza/go-dockerclient"

// CreateNetwork - creates a docker network
func (c Client) CreateNetwork(name string, driver string) (*docker.Network, error) {
	network, err := c.dockerClient.CreateNetwork(docker.CreateNetworkOptions{
		Name:   name,
		Driver: driver,
	})
	if err != nil {
		return nil, err
	}

	return network, nil
}

// RemoveNetwork - removes a docker network
func (c Client) RemoveNetwork(network *docker.Network) error {
	if err := c.dockerClient.RemoveNetwork(network.ID); err != nil {
		return err
	}

	return nil
}

// Connect - connects a container to a network
func (c Client) Connect(network *docker.Network, container *docker.Container) error {
	return c.dockerClient.ConnectNetwork(network.ID, docker.NetworkConnectionOptions{
		Container: container.ID,
	})
}
