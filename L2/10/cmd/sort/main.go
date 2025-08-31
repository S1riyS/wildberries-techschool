package main

import (
	"os"

	"github.com/S1riyS/wildberries-techschool/L2/10/internal/config"
	"github.com/S1riyS/wildberries-techschool/L2/10/internal/externalsorter"
	"github.com/S1riyS/wildberries-techschool/L2/10/internal/parser"
)

func main() {
	cfg := config.New()
	parser := parser.New(cfg)

	sorter := externalsorter.MustNew(cfg, parser)
	defer sorter.Cleanup()

	err := sorter.Sort(os.Stdin, os.Stdout)
	if err != nil {
		panic(err)
	}
}
