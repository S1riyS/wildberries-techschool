package main

import (
	"io"
	"os"

	"github.com/S1riyS/wildberries-techschool/L2/10/internal/config"
	"github.com/S1riyS/wildberries-techschool/L2/10/internal/externalsorter"
	"github.com/S1riyS/wildberries-techschool/L2/10/internal/parser"
)

func main() {
	cfg := config.New()
	parser := parser.New(cfg)

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

	// Init sorter
	sorter := externalsorter.MustNew(cfg, parser)

	// Check if sorted
	if cfg.IsCheckSorted {
		isSorted, err := sorter.CheckIfSorted(source, destination)
		if err != nil {
			panic(err)
		}
		if !isSorted {
			//nolint:gocritic // Cleanup is required after Sort, therefore it is safe to use os.Exit(1) before Sort
			os.Exit(1)
		}
		os.Exit(0) // Explicitly exit with 0
	}

	// Sort
	err := sorter.Sort(source, destination)
	defer sorter.Cleanup()
	if err != nil {
		panic(err)
	}
}
