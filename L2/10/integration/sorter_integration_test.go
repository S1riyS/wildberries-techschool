package integration_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	mainFile  = "../cmd/sort/main.go"
	buildFile = "../build/sort"
)

// TestMain is a special function that is called before any tests are run.
// It is used to build the binary before running the tests.
func TestMain(m *testing.M) {
	// Build the binary before running the tests
	if err := buildBinary(); err != nil {
		fmt.Printf("Failed to build binary: %v\n", err)
		os.Exit(1)
	}

	// Run the tests
	exitCode := m.Run()

	// Clean up after the tests have finished
	os.Remove(buildFile)
	os.Exit(exitCode)
}

// buildBinary builds the binary for the sorter.
func buildBinary() error {
	cmd := exec.Command("go", "build", "-o", buildFile, mainFile)
	return cmd.Run()
}

// runSorterWithStdin runs the sorter with the given arguments and input.
func runSorterWithStdin(args []string, input string) (string, string, error) {
	cmd := exec.Command(buildFile, args...)

	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// runSorterWithFileInput runs the sorter with the given arguments and input file.
func runSorterWithFileInput(args []string, inputFile string) (string, string, error) {
	cmd := exec.Command(buildFile, append(args, inputFile)...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// readTestFile reads the contents of the given file.
func readTestFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// TestSortWithTestData tests the sorter with test data.
func TestSortWithTestData(t *testing.T) {
	testCases := []struct {
		name    string
		testDir string
		args    []string
	}{
		{
			name:    "basic sort",
			testDir: "basic",
			args:    []string{},
		},
		{
			name:    "reverse sort",
			testDir: "reverse",
			args:    []string{"-r"},
		},
		{
			name:    "numeric sort",
			testDir: "numeric",
			args:    []string{"-n"},
		},
		{
			name:    "unique sort",
			testDir: "unique",
			args:    []string{"-u"},
		},
		{
			name:    "combined options",
			testDir: "combined",
			args:    []string{"-n", "-ru"},
		},
		{
			name:    "empty input",
			testDir: "empty",
			args:    []string{},
		},
		{
			name:    "check if sorted",
			testDir: "sorted",
			args:    []string{"-c"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testDir := filepath.Join("testdata", tc.testDir)

			// Read the input data
			inputData, err := readTestFile(filepath.Join(testDir, "input.txt"))
			if err != nil {
				t.Fatalf("Failed to read input file: %v", err)
			}

			// Run the sorter
			stdout, stderr, err := runSorterWithStdin(tc.args, inputData)

			// For successful tests, check the stdout
			if err != nil {
				t.Fatalf("Unexpected error: %v\nStderr: %s", err, stderr)
			}

			// Read the expected output
			expectedOutput, err := readTestFile(filepath.Join(testDir, "expected.txt"))
			if err != nil {
				t.Fatalf("Failed to read expected output file: %v", err)
			}

			// Normalize the output
			actual := normalizeLineEndings(stdout)
			expected := normalizeLineEndings(expectedOutput)

			if actual != expected {
				t.Errorf("Output mismatch.\nExpected:\n%s\nGot:\n%s",
					formatForDisplay(expected),
					formatForDisplay(actual))
			}
		})
	}
}

// TestSortWithFileInput tests the sorter with file input.
func TestSortWithFileInput(t *testing.T) {
	testCases := []struct {
		name    string
		testDir string
		args    []string
	}{
		{
			name:    "file input basic sort",
			testDir: "basic",
			args:    []string{},
		},
		{
			name:    "file input numeric sort",
			testDir: "numeric",
			args:    []string{"-n"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testDir := filepath.Join("testdata", tc.testDir)
			inputFile := filepath.Join(testDir, "input.txt")

			// Run the sorter with file input
			stdout, stderr, err := runSorterWithFileInput(tc.args, inputFile)
			if err != nil {
				t.Fatalf("Unexpected error: %v\nStderr: %s", err, stderr)
			}

			// Read the expected output
			expectedOutput, err := readTestFile(filepath.Join(testDir, "expected.txt"))
			if err != nil {
				t.Fatalf("Failed to read expected output file: %v", err)
			}

			// Normalize the output
			actual := normalizeLineEndings(stdout)
			expected := normalizeLineEndings(expectedOutput)

			if actual != expected {
				t.Errorf("Output mismatch.\nExpected:\n%s\nGot:\n%s",
					formatForDisplay(expected),
					formatForDisplay(actual))
			}
		})
	}
}

// normalizeLineEndings normalizes the line endings in the given string.
func normalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// formatForDisplay formats the given string for display in a test failure message.
func formatForDisplay(s string) string {
	return strings.ReplaceAll(s, "\n", "\\n\n")
}
