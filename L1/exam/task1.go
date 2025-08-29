package main

import "fmt"

func main() {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	s := a[2:5] // s = [3, 4, 5] len=3 cap=8

	s2 := s[5:7] // s2 = [8, 9] len=2 cap=3

	s2 = append(s2, 11) // s2 = [8, 9, 11] len=3 cap=3
	s2 = append(s2, 12) // s2 = [8, 9, 11, 12] len=4 cap=6

	fmt.Println(a) // a = [1, 2, 3, 4, 5, 6, 7, 8, 9, 11] len=10 cap=10
}
