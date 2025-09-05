package main

import (
	"fmt"
	"io"
	"os"

	"github.com/S1riyS/wildberries-techschool/L2/13/internal/config"
	"github.com/S1riyS/wildberries-techschool/L2/13/internal/cut"
)

func main() {
	exitCode := mainWithExitCode()
	os.Exit(exitCode)
}

func mainWithExitCode() int {
	cfg, err := config.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// Determine source
	var source io.Reader
	if cfg.InputFile != "" {
		var file *os.File
		file, err = os.Open(cfg.InputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file '%s': %v\n", cfg.InputFile, err)
			return 1
		}
		defer file.Close()
		source = file
	} else {
		source = os.Stdin
	}

	// Determine destination
	var destination io.Writer = os.Stdout

	// Init and run command
	c := cut.New(cfg)
	err = c.Run(source, destination)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	return 0
}
