package util

import (
	"fmt"
	"os/exec"
	"sync"

	log "github.com/Sirupsen/logrus"
	uuid "github.com/nu7hatch/gouuid"
)

// PostDockerActionDelayInSeconds - after every docker action we must sleep to wait for docker's caches to flush
const PostDockerActionDelayInSeconds = 2

// GetPrefixedUUID - returns a prefixed uuid
func GetPrefixedUUID(prefix string) (string, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s", prefix, u4), nil
}

// RunCommand - runs a shell command
func RunCommand(cmd string) ([]byte, error) {
	log.WithFields(log.Fields{"cmd": cmd}).Debug("Running cmd")
	return exec.Command("sh", "-c", cmd).Output()
}

// Work - executes a closure concurrently
func Work(workerCount int, worker func(), postWork func()) {
	wg := &sync.WaitGroup{}
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			worker()
		}()
	}

	go func() {
		wg.Wait()
		postWork()
	}()
}
