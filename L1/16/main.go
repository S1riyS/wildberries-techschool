package main

import "fmt"

// Average Time complexity: O(n * log(n)).
// Worst Time complexity: O(n^2).
func quickSort(arr []int) []int {
	// Base case
	if len(arr) <= 1 {
		return arr
	}

	// Choose pivot element (middle element)
	pivot := arr[len(arr)/2]

	// Distribute elements into sub-arrays
	var less, equal, greater []int
	for _, num := range arr {
		switch {
		case num < pivot:
			less = append(less, num)
		case num == pivot:
			equal = append(equal, num)
		case num > pivot:
			greater = append(greater, num)
		}
	}

	// Run recursively
	left := quickSort(less)
	right := quickSort(greater)
	return append(append(left, equal...), right...)
}

func main() {
	arr := []int{-5, 2, 9, 1, 6, 3, 7, 4, 8, 2}
	sortedArr := quickSort(arr)
	fmt.Println(sortedArr)
}
