package main

import (
	"fmt"

	"github.com/S1riyS/wildberries-techschool/L2/9/pkg/unpacker"
)

// ! NOTE: main.go is for manual testing.
// Implementation and tests are in pkg/unpacker directory.
func main() {
	data := `a0qwe\45`
	result, err := unpacker.Unpack(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}
