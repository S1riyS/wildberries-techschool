package main

import (
	"fmt"
	"io"
	"os"

	"github.com/S1riyS/wildberries-techschool/L2/13/internal/config"
)

func main() {
	exitCode := mainWithExitCode()
	os.Exit(exitCode)
}

// mainWithExitCode is an "internal" main function that returns an exit code.
// Purpose of this function is to allow defers to be used properly (os.Exit prevents defers from running).
func mainWithExitCode() int {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)

	// Determine source
	var source io.Reader
	if cfg.InputFile != "" {
		file, err := os.Open(cfg.InputFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		source = file
	} else {
		source = os.Stdin
	}
	_ = source

	// Determine destination
	var destination io.Writer = os.Stdout
	_ = destination

	// Init and run command
	// ...

	return 0
}
