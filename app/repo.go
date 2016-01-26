package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"os"
	"strings"
	"time"
)

type testOutputLine struct {
	Message  string `json:"message"`
	Line     string `json:"line"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
}

type testOutput struct {
	Output []testOutputLine `json:"output"`
}

// repo manager
type repoManager struct {
	repos            []repo
	cloneDestination string
}

func (r repoManager) cleanup() error {
	log.Info("Manager cleanup")
	return os.RemoveAll(r.cloneDestination)
}

// repo
type repo struct {
	manager       repoManager
	name          string
	testContainer *docker.Container
	client        *docker.Client
	testFailed    bool
}

func (r repo) clonePath() string {
	return fmt.Sprintf("%s/%s", r.manager.cloneDestination, r.name)
}

func (r repo) runRepoCommand(name string) ([]byte, error) {
	return runCommand(fmt.Sprintf("cd %s && %s", r.clonePath(), name))
}

func (r repo) log(fields map[string]interface{}) *log.Entry {
	fields["repo"] = r.name
	return log.WithFields(fields)
}

func (r repo) logInfo(message string) {
	r.log(map[string]interface{}{}).Info(message)
}

func (r repo) logWarning(message string) {
	r.log(map[string]interface{}{}).Warning(message)
}

func (r repo) runTests() (err error) {
	// cloning the repo in
	r.logInfo("Cloning")
	cloneCommand := fmt.Sprintf(
		"git clone %s %s",
		fmt.Sprintf("https://github.com/ihsw/%s.git", r.name),
		r.clonePath(),
	)
	if _, err = runCommand(cloneCommand); err != nil {
		r.logWarning("Could not clone")
		return r.cleanup(err)
	}

	// building up the related images
	r.logInfo("Building images")
	if _, err = r.runRepoCommand("./bin/build-images"); err != nil {
		r.logWarning("Could not build images")
		return r.cleanup(err)
	}

	// starting up the web-test service
	r.logInfo("Up web-test")
	if _, err = r.runRepoCommand("docker-compose up -d web-test"); err != nil {
		r.logWarning("Could not up web-tests")
		return r.cleanup(err)
	}

	// sleep so that the web-test container has time to get up
	time.Sleep(10 * time.Second)

	// fetching the name web-test container
	r.logInfo("Fetch web-test container name")
	var output []byte
	output, err = r.runRepoCommand("docker-compose ps -q web-test")
	if err != nil {
		r.logWarning("Could not fetch web-test container name")
		return r.cleanup(err)
	}
	webTestContainerName := strings.TrimSpace(string(output))

	// creating the test container
	r.logInfo("Create test container")
	r.testContainer, err = r.client.CreateContainer(docker.CreateContainerOptions{
		Name:   fmt.Sprintf("%s-tests", r.name),
		Config: &docker.Config{Image: "ihsw/the-matrix-tests"},
	})
	if err != nil {
		r.logWarning(fmt.Sprintf("Could not create test container: %s", err.Error()))
		return r.cleanup(err)
	}

	// starting the test container against the web-test container
	r.logInfo("Start test container against web-test container")
	err = r.client.StartContainer(r.testContainer.ID, &docker.HostConfig{
		Links: []string{fmt.Sprintf("%s:ApiServer", webTestContainerName)},
	})
	if err != nil {
		return r.cleanup(err)
	}

	// waiting for the test container to exit and checking the exit code
	r.logInfo("Waiting for the test container to exit")
	var status int
	if status, err = r.client.WaitContainer(r.testContainer.ID); err != nil {
		return r.cleanup(err)
	}
	if r.testFailed = status != 0; r.testFailed == true {
		r.logWarning("Test container exited with non-zero status")
		return r.cleanup(errors.New("Test container exited with non-zero status"))
	}

	return r.cleanup(nil)
}

func (r repo) cleanup(prevErr error) (err error) {
	r.logInfo("Cleaning up")

	// cleaning up the test container
	if r.testContainer != nil {
		// optionally fetching the test container's output for debugging
		if r.testFailed {
			// gathering up the logs
			var output bytes.Buffer
			err := r.client.Logs(docker.LogsOptions{
				Container:    r.testContainer.ID,
				Stdout:       true,
				OutputStream: &output,
			})
			if err != nil {
				return err
			}

			// reading the logs and checking for errors
			var parsedOutput testOutput
			if err := json.NewDecoder(strings.NewReader(output.String())).Decode(&parsedOutput); err != nil {
				return err
			}

			for _, outputLine := range parsedOutput.Output {
				r.log(map[string]interface{}{
					"expected": outputLine.Expected,
					"actual":   outputLine.Actual,
				}).Warning(outputLine.Message)
			}
		}

		r.logInfo("Removing test container")
		err = r.client.RemoveContainer(docker.RemoveContainerOptions{
			ID:            r.testContainer.ID,
			RemoveVolumes: true,
		})
		if err != nil {
			return err
		}
	}

	// stopping the services
	r.logInfo("Stop all")
	if _, err = r.runRepoCommand("docker-compose stop"); err != nil {
		r.logInfo("Could not stop all")
		return err
	}

	if _, err = r.runRepoCommand("docker rm -v $(docker-compose ps -q)"); err != nil {
		r.logInfo("Could not remove containers")
		return err
	}

	// removing the cloned repo
	r.logInfo("Remove repo files")
	if err = os.RemoveAll(r.clonePath()); err != nil {
		r.logInfo("Could not remove repo files")
		return err
	}

	return prevErr
}
