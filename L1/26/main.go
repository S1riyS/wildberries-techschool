package main

import (
	"fmt"
	"strings"
)

// Time complexity: O(n)
// Memory complexity: O(n)
func checkIfUnique(s string) bool {
	s = strings.ToLower(s)
	seen := make(map[rune]struct{})
	for _, c := range s {
		if _, ok := seen[c]; ok {
			return false
		}
		seen[c] = struct{}{}
	}
	return true
}

func main() {
	fmt.Println(checkIfUnique("abcd"))      // true
	fmt.Println(checkIfUnique("abCdefAaf")) // false
	fmt.Println(checkIfUnique("aabcd"))     // false
}
