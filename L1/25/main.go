package main

import (
	"fmt"
	"time"
)

func sleep(duration time.Duration) {
	ch := make(chan struct{})

	// Spin up a goroutine that will close the channel
	go func() {
		time.AfterFunc(duration, func() {
			close(ch)
		})
	}()

	// Wait for the channel to close (thus blocking the current goroutine)
	<-ch
}

func main() {
	fmt.Println("Начало работы")
	sleep(2 * time.Second)
	fmt.Println("Прошло 2 секунды")
}
