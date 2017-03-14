package util

import (
	"os/exec"
	"sync"

	log "github.com/Sirupsen/logrus"
)

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
