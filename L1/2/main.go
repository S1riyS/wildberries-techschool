package main

import (
	"fmt"
	"sync"
)

func main() {
	numbers := []int{2, 4, 6, 8, 10}
	results := make(chan int, len(numbers)) // Buffered channel for results

	var wg sync.WaitGroup
	for _, num := range numbers {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			square := n * n
			results <- square // Send result to channel
		}(num)
	}

	// Close channel after all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Read and print results
	for res := range results {
		fmt.Println(res)
	}
}
