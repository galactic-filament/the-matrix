package main

import (
	"fmt"
	"os/exec"
	"sync"
)

func runCommand(name string) error {
	return exec.Command("sh", "-c", name).Run()
}

func work(workerCount int, worker func(), postWork func()) {
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

func main() {
	fmt.Printf("Hello, world!")
}
