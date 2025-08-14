package main

import "fmt"

// Time complexity: O(n)
func makeUnique(arr []string) []string {
	presence := make(map[string]struct{})

	var result []string
	for _, v := range arr {
		if _, ok := presence[v]; ok {
			continue
		}
		presence[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

func main() {
	data := []string{"cat", "cat", "dog", "cat", "tree"}
	fmt.Println(makeUnique(data))
}
