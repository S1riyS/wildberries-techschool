package config

import "github.com/spf13/pflag"

type Config struct {
	// Flags
	RowsBefore       int
	RowsAfter        int
	RowsAround       int
	IsOnlyRowsCount  bool
	IsIgnoreRegister bool
	IsInverted       bool
	IsFixedRow       bool
	IsPrintRowNumber bool

	// Args
	Pattern   string
	InputFile string
}

func New() *Config {
	cfg := Config{}

	// Flags
	pflag.IntVarP(&cfg.RowsAfter, "after-content", "A", 0, "Print N lines of trailing context after matching lines.")
	pflag.IntVarP(&cfg.RowsBefore, "before-content", "B", 0, "Print N lines of trailing context before matching lines.")
	pflag.IntVarP(&cfg.RowsAround, "context", "C", 0, "Print N lines of context. Equivalent to -A N -B N.")
	pflag.BoolVarP(&cfg.IsOnlyRowsCount, "count", "c", false, "Only print the total number of matching lines.")
	pflag.BoolVarP(&cfg.IsIgnoreRegister, "ignore-case", "i", false, "Ignore case distinctions.")
	pflag.BoolVarP(&cfg.IsInverted, "invert", "v", false, "Select non-matching lines.")
	pflag.BoolVarP(
		&cfg.IsFixedRow, "fixed-strings", "F", false,
		"Interpret PATTERNS as fixed strings, not regular expressions.",
	)
	pflag.BoolVarP(&cfg.IsPrintRowNumber, "number", "n", false, "Number all output lines.")
	pflag.Parse()

	// Args
	cfg.Pattern = pflag.Arg(0)
	cfg.InputFile = pflag.Arg(1)

	return &cfg
}
