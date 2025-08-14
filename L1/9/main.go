package main

import "fmt"

func generator(data []int) chan int {
	outputCh := make(chan int)

	go func() {
		defer close(outputCh)
		for _, v := range data {
			outputCh <- v
		}
	}()

	return outputCh
}

func processor(inputCh chan int, f func(int) int) chan int {
	outputCh := make(chan int)

	go func() {
		defer close(outputCh)
		for v := range inputCh {
			outputCh <- f(v)
		}
	}()

	return outputCh
}

func main() {
	x := []int{0, 2, 4, 6, 8, 10}
	inputCh := generator(x)
	resultCh := processor(inputCh, func(v int) int {
		return v * 2
	})

	for v := range resultCh {
		fmt.Println(v)
	}
}
