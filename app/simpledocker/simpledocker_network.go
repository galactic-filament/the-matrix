package simpledocker

import docker "github.com/fsouza/go-dockerclient"

import "errors"
import "net"

// CreateNetwork - creates a docker network
func (c Client) CreateNetwork(name string, driver string) (*docker.Network, error) {
	network, err := c.dockerClient.CreateNetwork(docker.CreateNetworkOptions{
		Name:           name,
		Driver:         driver,
		CheckDuplicate: true,
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

// GetNetwork - fetches a docker network
func (c Client) GetNetwork(id string) (*docker.Network, error) { return c.dockerClient.NetworkInfo(id) }

// Connect - connects a container to a network
func (c Client) Connect(network *docker.Network, container *docker.Container) (*docker.Network, error) {
	err := c.dockerClient.ConnectNetwork(network.ID, docker.NetworkConnectionOptions{
		Container: container.ID,
	})
	if err != nil {
		return nil, err
	}

	return c.GetNetwork(network.ID)
}

// GetContainerIP - fetchs a container ip address from that network
func (c Client) GetContainerIP(network *docker.Network, container *docker.Container) (net.IP, error) {
	if network == nil {
		return nil, errors.New("Network was nil")
	}
	if container == nil {
		return nil, errors.New("Container was nil")
	}

	currentNetwork, err := c.GetNetwork(network.ID)
	if err != nil {
		return nil, err
	}

	endpoint, ok := currentNetwork.Containers[container.ID]
	if !ok {
		return nil, errors.New("Container was not a member of network")
	}

	ip, _, err := net.ParseCIDR(endpoint.IPv4Address)
	if err != nil {
		return nil, err
	}

	return ip, nil
}
