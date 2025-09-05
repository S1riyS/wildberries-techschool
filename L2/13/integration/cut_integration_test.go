package integration_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/S1riyS/wildberries-techschool/L2/13/internal/config"
	"github.com/S1riyS/wildberries-techschool/L2/13/internal/cut"
)

func TestCutIntegrationTestdata(t *testing.T) {
	testCases := []struct {
		name    string
		testDir string
		config  *config.Config
	}{
		{
			name:    "basic field extraction",
			testDir: "basic",
			config: &config.Config{
				Fields:    config.MustNewFields([]string{"1", "3"}),
				Delimiter: "\t",
			},
		},
		{
			name:    "only delimited lines",
			testDir: "only_delimited",
			config: &config.Config{
				Fields:          config.MustNewFields([]string{"1", "2"}),
				Delimiter:       ":",
				IsOnlyDelimited: true,
			},
		},
		{
			name:    "multiple fields",
			testDir: "multiple_fields",
			config: &config.Config{
				Fields:    config.MustNewFields([]string{"1-3", "5"}),
				Delimiter: " ",
			},
		},
		{
			name:    "range of fields",
			testDir: "range_fields",
			config: &config.Config{
				Fields:    config.MustNewFields([]string{"2-4"}),
				Delimiter: "\t",
			},
		},
		{
			name:    "mixed fields and ranges",
			testDir: "mixed_fields",
			config: &config.Config{
				Fields:    config.MustNewFields([]string{"1", "3-5", "7"}),
				Delimiter: ",",
			},
		},
		{
			name:    "custom delimiter",
			testDir: "custom_delimiter",
			config: &config.Config{
				Fields:    config.MustNewFields([]string{"2", "3"}),
				Delimiter: ";",
			},
		},
		{
			name:    "no delimiter found",
			testDir: "no_delimiter",
			config: &config.Config{
				Fields:          config.MustNewFields([]string{"1", "2"}),
				Delimiter:       ":",
				IsOnlyDelimited: true,
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

			// Run cut
			var output bytes.Buffer
			c := cut.New(tc.config)
			err = c.Run(inputFile, &output)
			if err != nil {
				t.Fatalf("Cut failed: %v", err)
			}

			// Compare results
			result := output.String()
			if result != expected {
				t.Errorf("Test %s failed.\nExpected:\n%q\nGot:\n%q", tc.name, expected, result)
			}
		})
	}
}
