package main

import (
	"io"
	"os"

	"github.com/S1riyS/wildberries-techschool/L2/10/internal/config"
	"github.com/S1riyS/wildberries-techschool/L2/10/internal/grep"
)

func main() {
	exitCode := mainWithExitCode()
	os.Exit(exitCode)
}

// mainWithExitCode is an "internal" main function that returns an exit code.
// Purpose of this function is to allow defers to be used properly (os.Exit prevents defers from running).
func mainWithExitCode() int {
	cfg := config.New()
	printer := grep.NewPrinter(cfg)

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

	// Determine destination
	var destination io.Writer = os.Stdout

	// Init and run command
	g := grep.New(cfg, printer)
	matchCount := g.Run(source, destination)

	// Return exit code
	if matchCount == 0 {
		return 1
	}
	return 0
}
