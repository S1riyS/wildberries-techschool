package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func WorkerWithChan(stopCh chan struct{}) {
	fmt.Println("workerWithChan started")
	for {
		select {
		case <-stopCh:
			fmt.Println("workerWithChan stopped")
			return
		default:
			fmt.Println("working...")
			time.Sleep(time.Second)
		}
	}
}

func WorkerWithFlag(stopFlag *bool) {
	fmt.Println("workerWithFlag started")
	for !*stopFlag {
		fmt.Println("working...")
		time.Sleep(time.Second)
	}
	fmt.Println("workerWithFlag stopped")
}

func WorkerWithContextDone(ctx context.Context) {
	fmt.Println("workerWithContextDone started")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("workerWithContextDone stopped")
			return
		default:
			fmt.Println("working...")
			time.Sleep(time.Second)
		}
	}
}

func WorkerWithContextTimeout(ctx context.Context) {
	fmt.Println("workerWithContextTimeout started")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("workerWithContextTimeout stopped")
			return
		default:
			fmt.Println("working...")
			time.Sleep(time.Second)
		}
	}
}

func WorkerWithGoexit(timeToWork time.Duration) {
	fmt.Println("workerWithGoexit started")
	for {
		fmt.Println("working...")
		time.Sleep(timeToWork)

		fmt.Println("workerWithGoexit stopped")
		runtime.Goexit()
	}
}

func main() {
	// valeiables
	N := 3
	TimeToWork := time.Duration(N) * time.Second

	// Goroutine with chan
	stopCh := make(chan struct{})
	go WorkerWithChan(stopCh)
	time.Sleep(TimeToWork)
	stopCh <- struct{}{}
	close(stopCh)

	// Goroutine with flag
	stopFlag := false
	go WorkerWithFlag(&stopFlag)
	time.Sleep(TimeToWork)
	stopFlag = true

	// Goroutine with context
	ctx, cancel := context.WithCancel(context.Background())
	go WorkerWithContextDone(ctx)
	time.Sleep(TimeToWork)
	cancel()

	// Goroutine with context with timeout
	ctx, cancel = context.WithTimeout(context.Background(), TimeToWork)
	go WorkerWithContextTimeout(ctx)
	defer cancel()

	// Goroutine with Goexit
	go WorkerWithGoexit(TimeToWork)

	// Handle shutdown
	// ! NOTE: could be done with WaitGroup, but I don't want to complicate the example
	time.Sleep(2 * TimeToWork)
	fmt.Println("Done")
}
