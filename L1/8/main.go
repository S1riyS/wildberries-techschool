package main

import "fmt"

func setBit(num int64, i uint, value uint) int64 {
	if value == 1 {
		return num | (1 << i)
	} else {
		return num &^ (1 << i)
	}
}

func main() {
	var num int64 = 5
	var index uint = 0
	var value uint = 0

	fmt.Printf("Before\t(bin):\t%b\n", num)
	num = setBit(num, index, value)
	fmt.Printf("After\t(bin):\t%b\n", num)
	fmt.Printf("After\t(dec):\t%d\n", num)
}
