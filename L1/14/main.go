package main

import (
	"fmt"
	"reflect"
)

func detectType(value any) string {
	switch value.(type) {
	case int:
		return "int"
	case string:
		return "string"
	case bool:
		return "bool"
	default:
		// As it turns out, you can not use type assertion for generic chans.
		// That's why I am using reflect package
		if reflect.TypeOf(value).Kind() == reflect.Chan {
			return "chan"
		}
		return "unknown"
	}
}

func main() {
	var (
		num    int    = 42
		str    string = "hello"
		flag   bool   = true
		chInt  chan int
		chStr  chan string
		random any = 3.14 // Generic type
	)

	chInt = make(chan int)
	chStr = make(chan string)

	fmt.Println("num is:", detectType(num))       // int
	fmt.Println("str is:", detectType(str))       // string
	fmt.Println("flag is:", detectType(flag))     // bool
	fmt.Println("chInt is:", detectType(chInt))   // chan
	fmt.Println("chStr is:", detectType(chStr))   // chan
	fmt.Println("random is:", detectType(random)) // unknown
}
