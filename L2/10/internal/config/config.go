package config

import "github.com/spf13/pflag"

type Config struct {
	// Flags
	Column         int
	IsNumeric      bool
	IsReverse      bool
	IsUnique       bool
	IsMonthSort    bool
	IsIgnoreBlanks bool
	IsCheckSorted  bool
	IsHumanNumeric bool

	// Args
	InputFile string
}

func New() *Config {
	f := Config{}

	// Flags
	pflag.IntVarP(&f.Column, "key", "k", 0, "sort by column N")
	pflag.BoolVarP(&f.IsNumeric, "numeric", "n", false, "sort numerically")
	pflag.BoolVarP(&f.IsReverse, "reverse", "r", false, "reverse sort")
	pflag.BoolVarP(&f.IsUnique, "unique", "u", false, "output only unique lines")
	pflag.BoolVarP(&f.IsMonthSort, "month-sort", "M", false, "sort by month, use with -k")
	pflag.BoolVarP(&f.IsIgnoreBlanks, "ignore-blanks", "b", false, "ignore trailing blanks")
	pflag.BoolVarP(&f.IsCheckSorted, "check", "c", false, "check if sorted")
	pflag.BoolVarP(&f.IsHumanNumeric, "human-numeric", "h", false, "sort human readable numbers")
	pflag.Parse()

	// Args
	f.InputFile = pflag.Arg(0)

	return &f
}
