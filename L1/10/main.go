package main

import "fmt"

func solve(data []float64) map[int][]float64 {
	groups := make(map[int][]float64)
	var key int
	for _, v := range data {
		key = int(v/10) * 10
		groups[key] = append(groups[key], v)
	}

	return groups
}

func main() {
	tempSequence := []float64{-35.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	groups := solve(tempSequence)
	for k, v := range groups {
		fmt.Printf("%v: %v\n", k, v)
	}
}
