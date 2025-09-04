package config

import (
	"github.com/spf13/pflag"
)

type Config struct {
	// Flags
	FieldsRanges    []*FieldsRange
	Delimeter       string
	IsOnlyDelimited bool

	// Args
	InputFile string
}

func New() (*Config, error) {
	cfg := Config{}

	// Flagss
	var fieldsRaw []string
	pflag.StringSliceVarP(&fieldsRaw, "fields", "f", nil, "Comma-separated list of fields. Example: --fields 1,3-5")
	pflag.StringVarP(&cfg.Delimeter, "delimiter", "d", "\t", "Delimiter symbol.")
	pflag.BoolVarP(&cfg.IsOnlyDelimited, "only-delimited", "s", false, "Print only lines with delimiters.")

	// Parsing
	pflag.Parse()
	if err := cfg.parseFields(fieldsRaw); err != nil {
		return nil, err
	}

	// Args
	cfg.InputFile = pflag.Arg(0)

	return &cfg, nil
}

func (c *Config) parseFields(fieldsRangesRaw []string) error {
	c.FieldsRanges = make([]*FieldsRange, len(fieldsRangesRaw))
	for i := range len(fieldsRangesRaw) {
		fieldsRange, err := NewFieldsRange(fieldsRangesRaw[i])
		if err != nil {
			return err
		}
		c.FieldsRanges[i] = fieldsRange
	}
	return nil
}
