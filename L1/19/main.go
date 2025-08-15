package main

import "fmt"

// Time complexit: O(n).
// Memory complexity: O(n)
func reverseString(value string) string {
	runes := []rune(value)
	n := len(runes)
	for i := range n / 2 {
		runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
	}
	return string(runes)
}

func main() {
	data := "Hello, World! 👋 Привет, мир!"
	fmt.Println(reverseString(data))
}
