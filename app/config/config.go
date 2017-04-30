package config

import (
	"encoding/json"
)

// Config - parsed config from a json string
type Config struct {
	DockerHost string `json:"docker_host"`
}

// Parse - parses a json blob for config
func Parse(data []byte) (Config, error) {
	var config Config
	err := json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
