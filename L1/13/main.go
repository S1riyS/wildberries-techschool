package main

import "fmt"

// Memory complexity: O(1)
func main() {
	a := 18
	b := 12

	// I wrote this thing in Assembler last year as a part of my university assignment. LOL
	fmt.Println(a, b)
	a ^= b
	b ^= a
	a ^= b
	fmt.Println(a, b)
}
