package anagrams_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/S1riyS/wildberries-techschool/L2/11/pkg/anagrams"
)

func TestFindAnagrams(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected map[string][]string
	}{
		{
			name:  "russian anagrams",
			input: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"},
			expected: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{
			name:     "no anagrams",
			input:    []string{"cat", "dog", "bird"},
			expected: map[string][]string{},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: map[string][]string{},
		},
		{
			name:     "single word",
			input:    []string{"hello"},
			expected: map[string][]string{},
		},
		{
			name:  "lower and upper cases",
			input: []string{"Listen", "SILENT", "enlist", "Google", "gooGLe"},
			expected: map[string][]string{
				"listen": {"enlist", "listen", "silent"},
				"google": {"google", "google"},
			},
		},
		{
			name:  "multiple anagram groups",
			input: []string{"abc", "bac", "cab", "cba", "def", "fed", "xyz"},
			expected: map[string][]string{
				"abc": {"abc", "bac", "cab", "cba"},
				"def": {"def", "fed"},
			},
		},
		{
			name:  "russian words with different cases",
			input: []string{"Кот", "ТОК", "кто", "отк", "Актер", "терка", "катер"},
			expected: map[string][]string{
				"кот":   {"кот", "кто", "отк", "ток"},
				"актер": {"актер", "катер", "терка"},
			},
		},
		{
			name:  "duplicate words",
			input: []string{"пятак", "пятак", "пятка", "пятка", "тяпка"},
			expected: map[string][]string{
				"пятак": {"пятак", "пятак", "пятка", "пятка", "тяпка"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := anagrams.Find(tt.input)

			// Normalize both expected and result for comparison
			normalizedExpected := normalizeAnagramMap(tt.expected)
			normalizedResult := normalizeAnagramMap(result)

			if !reflect.DeepEqual(normalizedExpected, normalizedResult) {
				t.Errorf("FindAnagrams(%v) =\n%v\nExpected:\n%v", tt.input, normalizedResult, normalizedExpected)
			}
		})
	}
}

// normalizeAnagramMap sorts all value slices in the map for consistent comparison.
func normalizeAnagramMap(m map[string][]string) map[string][]string {
	normalized := make(map[string][]string)
	for k, v := range m {
		// Create a copy and sort
		sortedSlice := make([]string, len(v))
		copy(sortedSlice, v)
		sort.Strings(sortedSlice)
		normalized[k] = sortedSlice
	}
	return normalized
}
