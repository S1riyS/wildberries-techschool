package main

import "fmt"

func deleteElement(slice []int, index int) []int {
	copy(slice[index:], slice[index+1:])
	return slice[:len(slice)-1]
}

func main() {
	a := []int{1, 2, 3, 4, 5}
	fmt.Println(deleteElement(a, 2))
}
