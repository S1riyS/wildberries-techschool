package config

import (
	"math"
	"sort"
)

type Fields struct {
	ranges   []*Range
	MinIndex int
	MaxIndex int
}

func NewFields(ranges []string) (*Fields, error) {
	f := &Fields{
		ranges:   make([]*Range, len(ranges)),
		MinIndex: math.MaxInt,
		MaxIndex: math.MinInt,
	}

	for i, rangeString := range ranges {
		// Parse range
		r, err := NewRange(rangeString)
		if err != nil {
			return nil, err
		}
		// Save range
		f.ranges[i] = r
		f.MinIndex = min(f.MinIndex, r.From)
		f.MaxIndex = max(f.MaxIndex, r.To)
	}

	// Sort ranges by start index
	// * Time complexity: O(n * log(n))
	sort.Slice(f.ranges, func(i, j int) bool {
		return f.ranges[i].From < f.ranges[j].From
	})

	// Merge overlapping ranges
	// * Time complexity: O(n)
	for i := 0; i < len(f.ranges)-1; i++ {
		if f.ranges[i].To >= f.ranges[i+1].From-1 {
			f.ranges[i].To = max(f.ranges[i].To, f.ranges[i+1].To)
			f.ranges = append(f.ranges[:i+1], f.ranges[i+2:]...)
			i--
		}
	}

	return f, nil
}

// MustNewFields is a wrapper around NewFields that panics if an error occurs.
// Should only be used in tests.
func MustNewFields(ranges []string) *Fields {
	f, err := NewFields(ranges)
	if err != nil {
		panic(err)
	}
	return f
}

func (f *Fields) IsInRange(index int) bool {
	if index < f.MinIndex || index > f.MaxIndex {
		return false
	}

	left := 0
	right := len(f.ranges) - 1

	for left <= right {
		mid := left + (right-left)/2 //nolint:mnd // Default binary search implementation
		r := f.ranges[mid]

		if index >= r.From && index <= r.To {
			return true
		}

		if index < r.From {
			// Index is before current range, search left
			right = mid - 1
		} else {
			// Index is after current range, search right
			left = mid + 1
		}
	}

	return false
}
