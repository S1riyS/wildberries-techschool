package config

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidFieldRange = errors.New("invalid field range")
)

const (
	correctRangeArgsCount = 2
	defaultRangeDelimiter = "-"
)

type Range struct {
	From int
	To   int
}

func NewRange(rangeString string) (*Range, error) {
	if strings.Contains(rangeString, defaultRangeDelimiter) {
		// Range containse a delimiter
		splitRange := strings.Split(rangeString, defaultRangeDelimiter)
		if len(splitRange) != correctRangeArgsCount {
			return nil, ErrInvalidFieldRange
		}

		from, err := strconv.Atoi(splitRange[0])
		if err != nil {
			return nil, ErrInvalidFieldRange
		}
		to, err := strconv.Atoi(splitRange[1])
		if err != nil {
			return nil, ErrInvalidFieldRange
		}

		if from > to {
			return nil, ErrInvalidFieldRange
		}

		return &Range{From: from, To: to}, nil
	}

	// Range doesn't contain a delimiter
	value, err := strconv.Atoi(rangeString)
	if err != nil {
		return nil, ErrInvalidFieldRange
	}
	return &Range{From: value, To: value}, nil
}
