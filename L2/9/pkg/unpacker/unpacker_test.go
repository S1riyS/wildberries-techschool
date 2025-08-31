package unpacker_test

import (
	"testing"

	"github.com/S1riyS/wildberries-techschool/L2/9/pkg/unpacker"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "basic unpacking",
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
			hasError: false,
		},
		{
			name:     "no digits",
			input:    "abcd",
			expected: "abcd",
			hasError: false,
		},
		{
			name:     "only digits",
			input:    "45",
			expected: "",
			hasError: true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			hasError: false,
		},
		{
			name:     "escaped digits 1",
			input:    "qwe\\4\\5",
			expected: "qwe45",
			hasError: false,
		},
		{
			name:     "escaped digits 2",
			input:    "qwe\\45",
			expected: "qwe44444",
			hasError: false,
		},
		{
			name:     "escaped backslash",
			input:    "qwe\\\\5",
			expected: "qwe\\\\\\\\\\",
			hasError: false,
		},
		{
			name:     "invalid escape sequence",
			input:    "qwe\\",
			expected: "",
			hasError: true,
		},
		{
			name:     "digit at start",
			input:    "0abc",
			expected: "",
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := unpacker.Unpack(test.input)

			if test.hasError {
				if err == nil {
					t.Errorf("Expected error for input %s, but got none", test.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %s: %v", test.input, err)
				}
				if result != test.expected {
					t.Errorf("For input %s, expected %s, but got %s", test.input, test.expected, result)
				}
			}
		})
	}
}
