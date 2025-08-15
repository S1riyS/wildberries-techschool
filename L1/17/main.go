package main

import "fmt"

func binarySearch(arr []int, target int) int {
	l := 0
	r := len(arr) - 1
	for r >= l {
		m := (l + r) / 2
		if arr[m] == target {
			return m
		} else if arr[m] > target {
			r = m - 1
		} else {
			l = m + 1
		}
	}
	return -1
}

func main() {
	data := []int{1, 2, 4, 6, 7, 8}
	target := 4
	fmt.Println(binarySearch(data, target))
}
