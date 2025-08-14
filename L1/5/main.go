package main

import (
	"fmt"
	"time"
)

func main() {
	N := 5
	timeToWork := time.Duration(N) * time.Second
	dataCh := make(chan int)

	// Send data to the channel
	go func() {
		i := 0
		for {
			i++
			dataCh <- i
			time.Sleep(300 * time.Millisecond)
		}
	}()

	// Receive data from the channel
	go func() {
		for value := range dataCh {
			fmt.Printf("Received value: %d\n", value)
		}
	}()

	// Handle shutdown
	<-time.After(timeToWork)
	close(dataCh)
	fmt.Println("Done")
}
