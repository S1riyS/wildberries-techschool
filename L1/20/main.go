package main

import (
	"fmt"
)

// Time complexity: O(n).
// Memory complexity: O(1).
func reverseWords(s string) string {
	runes := []rune(s)
	n := len(runes)

	// Reverse whole string:
	reverse(runes, 0, n-1)

	// Reverse each word
	start := 0
	for i := 0; i <= n; i++ {
		if i == n || runes[i] == ' ' {
			reverse(runes, start, i-1)
			start = i + 1
		}
	}

	return string(runes)
}

// Time complexity: O(n).
// Memory complexity: O(1).
func reverse(runes []rune, left, right int) {
	for left < right {
		runes[left], runes[right] = runes[right], runes[left]
		left++
		right--
	}
}

func main() {
	// ! NOTE: Solution idea:
	// (Reverse whole string) -> (Reverse each word) -> Success!
	// E.g.: "snow dog sun" -> "nus god wons" -> "sun dog snow"

	input := "snow dog sun"
	fmt.Println(reverseWords(input))
}
