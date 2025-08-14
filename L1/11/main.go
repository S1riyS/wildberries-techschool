package main

import "fmt"

// Find intersection of two arrays
//
// Time complexity: O(n + m)
func intersect(a, b []int) []int {
	counter := make(map[int]int)
	for _, v := range a {
		counter[v]++
	}

	var result []int
	for _, v := range b {
		if counter[v] > 0 {
			result = append(result, v)
			counter[v]--
		}
	}
	return result
}

func main() {
	a := []int{1, 2, 3, 3}
	b := []int{2, 3, 3, 4}
	fmt.Println(intersect(a, b))
}
