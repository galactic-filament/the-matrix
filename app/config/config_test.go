package config

import "testing"
import "fmt"

func TestParse(t *testing.T) {
	dockerHost := "tcp://localhost:2375"
	jsonBlob := []byte(fmt.Sprintf("{\"docker_host\": \"%s\"}", dockerHost))
	config, err := Parse(jsonBlob)
	if err != nil {
		t.Errorf("Could not parse json blob: %s", err.Error())
		return
	}

	if config.DockerHost != dockerHost {
		t.Error("Config DockerHost does not match given dockerHost")
		return
	}
}
