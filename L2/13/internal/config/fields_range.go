package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	defaultRangeDelimiter = "-"
)

// TODO: Combine all Ranges under one struct that tracks min and max of ranges (in order to iterate optimally)
// + It makes sence to keep ranges sorted by start or end in order to check if value is in any range with bin search

type FieldsRange struct {
	From int
	To   int
}

func (f *FieldsRange) String() string {
	return fmt.Sprintf("%d-%d", f.From, f.To)
}

func NewFieldsRange(rangeString string) (*FieldsRange, error) {
	if strings.Contains(rangeString, defaultRangeDelimiter) {
		// Range containse a delimiter
		splitRange := strings.Split(rangeString, defaultRangeDelimiter)
		if len(splitRange) != 2 {
			return nil, errors.New("invalid range format")
		}

		from, err := strconv.Atoi(splitRange[0])
		if err != nil {
			return nil, err
		}
		to, err := strconv.Atoi(splitRange[1])
		if err != nil {
			return nil, err
		}
		return &FieldsRange{From: from, To: to}, nil
	}

	// Range doesn't contain a delimiter
	value, err := strconv.Atoi(rangeString)
	if err != nil {
		return nil, err
	}
	return &FieldsRange{From: value, To: value}, nil
}

func (f *FieldsRange) IsIn(value int) bool {
	return value >= f.From && value <= f.To
}
