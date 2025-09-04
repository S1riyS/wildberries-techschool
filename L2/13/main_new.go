package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

func main() {
	// Define a string slice flag (comma-separated)
	var tags []string
	pflag.StringSliceVarP(&tags, "tags", "t", nil, "Comma-separated list of tags")

	pflag.Parse()

	fmt.Printf("Tags: %v\n", tags)
	fmt.Printf("Tags len: %d\n", len(tags))
}
