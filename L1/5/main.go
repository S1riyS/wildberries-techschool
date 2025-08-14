package main

import (
	"fmt"
	"time"
)

func main() {
	timeToWork := 5 * time.Second
	dataCh := make(chan int)

	// Send data to the channel
	go func() {
		i := 0
		for {
			i++
			dataCh <- i
			time.Sleep(300 * time.Millisecond) // небольшая задержка между отправками
		}
	}()

	// Receive data from the channel
	go func() {
		for value := range dataCh {
			fmt.Printf("Received value: %d\n", value)
		}
	}()

	// Ждем N секунд, затем завершаем программу
	<-time.After(timeToWork)
	close(dataCh)
	fmt.Println("Done")
}
