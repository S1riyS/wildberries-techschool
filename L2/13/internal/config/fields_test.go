package config_test

import (
	"errors"
	"math"
	"testing"

	"github.com/S1riyS/wildberries-techschool/L2/13/internal/config"
)

//nolint:gocognit // I don't want to refactor tests, please accept my apologies
func TestNewFields(t *testing.T) {
	tests := []struct {
		name    string
		ranges  []string
		wantMin int
		wantMax int
		wantErr bool
	}{
		{
			name:    "single range",
			ranges:  []string{"1-5"},
			wantMin: 1,
			wantMax: 5,
		},
		{
			name:    "multiple ranges",
			ranges:  []string{"1-5", "10-15", "20-25"},
			wantMin: 1,
			wantMax: 25,
		},
		{
			name:    "single number ranges",
			ranges:  []string{"1", "5", "10"},
			wantMin: 1,
			wantMax: 10,
		},
		{
			name:    "negative ranges",
			ranges:  []string{"-5--1", "1-5"},
			wantErr: true,
		},
		{
			name:    "overlapping ranges should merge",
			ranges:  []string{"1-5", "3-8", "10-15"},
			wantMin: 1,
			wantMax: 15,
		},
		{
			name:    "adjacent ranges should merge",
			ranges:  []string{"1-5", "6-10", "11-15"},
			wantMin: 1,
			wantMax: 15,
		},
		{
			name:    "invalid range format",
			ranges:  []string{"1-5", "invalid"},
			wantErr: true,
		},
		{
			name:    "empty range string",
			ranges:  []string{""},
			wantErr: true,
		},
		{
			name:    "empty ranges slice",
			ranges:  []string{},
			wantMin: math.MaxInt,
			wantMax: math.MinInt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := config.NewFields(tt.ranges)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewFields() expected error, got nil")
				}
				if !errors.Is(err, config.ErrInvalidFieldRange) {
					t.Errorf("NewFields() error = %v, want %v", err, config.ErrInvalidFieldRange)
				}

				return
			}

			if err != nil {
				t.Errorf("NewFields() unexpected error: %v", err)
				return
			}

			if fields.MinIndex != tt.wantMin {
				t.Errorf("NewFields() MinIndex = %d, want %d", fields.MinIndex, tt.wantMin)
			}

			if fields.MaxIndex != tt.wantMax {
				t.Errorf("NewFields() MaxIndex = %d, want %d", fields.MaxIndex, tt.wantMax)
			}
		})
	}
}

func TestFields_IsInRange(t *testing.T) {
	tests := []struct {
		name     string
		ranges   []string
		index    int
		expected bool
	}{
		{
			name:     "within single range",
			ranges:   []string{"1-5"},
			index:    3,
			expected: true,
		},
		{
			name:     "at start of range",
			ranges:   []string{"1-5"},
			index:    1,
			expected: true,
		},
		{
			name:     "at end of range",
			ranges:   []string{"1-5"},
			index:    5,
			expected: true,
		},
		{
			name:     "before range",
			ranges:   []string{"1-5"},
			index:    0,
			expected: false,
		},
		{
			name:     "after range",
			ranges:   []string{"1-5"},
			index:    6,
			expected: false,
		},
		{
			name:     "within first of multiple ranges",
			ranges:   []string{"1-5", "10-15"},
			index:    3,
			expected: true,
		},
		{
			name:     "within second of multiple ranges",
			ranges:   []string{"1-5", "10-15"},
			index:    12,
			expected: true,
		},
		{
			name:     "between ranges",
			ranges:   []string{"1-5", "10-15"},
			index:    7,
			expected: false,
		},
		{
			name:     "within merged overlapping ranges",
			ranges:   []string{"1-5", "3-8"},
			index:    6,
			expected: true,
		},
		{
			name:     "within merged adjacent ranges",
			ranges:   []string{"1-5", "6-10"},
			index:    7,
			expected: true,
		},
		{
			name:     "single number range match",
			ranges:   []string{"5"},
			index:    5,
			expected: true,
		},
		{
			name:     "single number range no match",
			ranges:   []string{"5"},
			index:    4,
			expected: false,
		},
		{
			name:     "empty ranges always false",
			ranges:   []string{},
			index:    1,
			expected: false,
		},
		{
			name:     "index at MinIndex boundary",
			ranges:   []string{"1-5", "10-15"},
			index:    1,
			expected: true,
		},
		{
			name:     "index at MaxIndex boundary",
			ranges:   []string{"1-5", "10-15"},
			index:    15,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := config.NewFields(tt.ranges)
			if err != nil {
				t.Fatalf("NewFields() failed: %v", err)
			}

			result := fields.IsInRange(tt.index)
			if result != tt.expected {
				t.Errorf("IsInRange(%d) = %v, want %v", tt.index, result, tt.expected)
			}
		})
	}
}

func TestFields_EdgeCases(t *testing.T) {
	t.Run("very large ranges", func(t *testing.T) {
		fields, err := config.NewFields([]string{"1-1000000"})
		if err != nil {
			t.Fatalf("NewFields() failed: %v", err)
		}

		if !fields.IsInRange(500000) {
			t.Error("IsInRange() should return true for index in large range")
		}

		if fields.IsInRange(1000001) {
			t.Error("IsInRange() should return false for index outside large range")
		}
	})

	t.Run("binary search correctness", func(t *testing.T) {
		// Test that binary search works correctly with many ranges
		ranges := []string{"1-5", "10-15", "20", "21", "22", "23-25", "30-35", "40-45"}
		fields, err := config.NewFields(ranges)
		if err != nil {
			t.Fatalf("NewFields() failed: %v", err)
		}

		// Test boundaries and midpoints
		testCases := []struct {
			index    int
			expected bool
		}{
			{3, true},   // middle of first range
			{13, true},  // middle of second range
			{23, true},  // middle of third range
			{33, true},  // middle of fourth range
			{43, true},  // middle of fifth range
			{8, false},  // between first and second
			{18, false}, // between second and third
			{28, false}, // between third and fourth
			{38, false}, // between fourth and fifth
			{0, false},  // before first
			{50, false}, // after last
		}

		for _, tc := range testCases {
			result := fields.IsInRange(tc.index)
			if result != tc.expected {
				t.Errorf("IsInRange(%d) = %v, want %v", tc.index, result, tc.expected)
			}
		}
	})
}
