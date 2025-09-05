package config

import (
	"github.com/spf13/pflag"
)

type Config struct {
	// Flags
	Fields          *Fields
	Delimiter       string
	IsOnlyDelimited bool

	// Args
	InputFile string
}

func New() (*Config, error) {
	cfg := Config{}

	// Flags
	var fieldsRaw []string
	pflag.StringSliceVarP(&fieldsRaw, "fields", "f", nil, "Comma-separated list of fields. Example: --fields 1,3-5")
	pflag.StringVarP(&cfg.Delimiter, "delimiter", "d", "\t", "Delimiter symbol.")
	pflag.BoolVarP(&cfg.IsOnlyDelimited, "only-delimited", "s", false, "Print only lines with delimiters.")

	// Parsing
	pflag.Parse()

	// Fields
	fields, err := NewFields(fieldsRaw)
	if err != nil {
		return nil, err
	}
	cfg.Fields = fields

	// Args
	cfg.InputFile = pflag.Arg(0)

	return &cfg, nil
}
