package integration_test

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/S1riyS/wildberries-techschool/L2/10/internal/config"
	"github.com/S1riyS/wildberries-techschool/L2/10/internal/grep"
)

func TestGrepIntegrationTestdata(t *testing.T) {
	testCases := []struct {
		name    string
		testDir string
		config  *config.Config
	}{
		{
			name:    "basic pattern matching",
			testDir: "basic",
			config: &config.Config{
				Pattern: "error",
			},
		},
		{
			name:    "case insensitive",
			testDir: "case_insensitive",
			config: &config.Config{
				Pattern:          "ErRoR",
				IsIgnoreRegister: true,
			},
		},
		{
			name:    "count only",
			testDir: "count_only",
			config: &config.Config{
				Pattern:         "test",
				IsOnlyRowsCount: true,
			},
		},
		{
			name:    "with context",
			testDir: "with_context",
			config: &config.Config{
				Pattern:    "target",
				RowsBefore: 1,
				RowsAfter:  1,
			},
		},
		{
			name:    "inverted match",
			testDir: "inverted",
			config: &config.Config{
				Pattern:    "skip",
				IsInverted: true,
			},
		},
		{
			name:    "fixed string matching",
			testDir: "fixed_string",
			config: &config.Config{
				Pattern:    "t.st",
				IsFixedRow: true,
			},
		},
		{
			name:    "regex matching",
			testDir: "regex",
			config: &config.Config{
				Pattern:    "t.st",
				IsFixedRow: false,
			},
		},
		{
			name:    "row number",
			testDir: "row_number",
			config: &config.Config{
				Pattern:          "test",
				IsPrintRowNumber: true,
				RowsAround:       1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Read input file
			inputPath := filepath.Join("testdata", tc.testDir, "input.txt")
			inputFile, err := os.Open(inputPath)
			if err != nil {
				t.Fatalf("Failed to open input file %s: %v", inputPath, err)
			}
			defer inputFile.Close()

			// Read expected output
			expectedPath := filepath.Join("testdata", tc.testDir, "expected.txt")
			expectedBytes, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatalf("Failed to read expected file %s: %v", expectedPath, err)
			}
			expected := string(expectedBytes)

			// Run grep
			var output bytes.Buffer
			g := grep.New(tc.config, grep.NewPrinter(tc.config))
			g.Run(inputFile, &output)

			// Clean up result
			result := removeANSICodes(output.String())

			// Compare results
			if result != expected {
				t.Errorf("Test %s failed.\nExpected:\n%q\nGot:\n%q", tc.name, expected, result)
			}
		})
	}
}

// removeANSICodes removes ANSI escape codes from the given string.
func removeANSICodes(s string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(s, "")
}
