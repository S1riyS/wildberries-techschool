package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Read `numWorkers` argument
	if len(os.Args) < 2 {
		fmt.Println("Please provide a numWorkers as a number argument")
		return
	}
	numWorkers, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid numWorkers: %s\n", os.Args[1])
		return
	}

	// Set up channels
	dataCh := make(chan int)
	stopCh := make(chan os.Signal, 1)

	// Spin up workers
	fmt.Printf("Spinning up %d workers\n", numWorkers)
	var wg sync.WaitGroup
	for i := range numWorkers {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for value := range dataCh {
				fmt.Printf("worker %d received value %d\n", workerID, value)
				time.Sleep(300 * time.Millisecond)
			}
		}(i)
	}

	// Handle shutdown
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)

	// Write data to channel
	for {
		select {
		case dataCh <- rand.Int():
		case <-stopCh:
			fmt.Println("Shutting down")
			close(dataCh)
			wg.Wait()
			return
		}
	}
}
