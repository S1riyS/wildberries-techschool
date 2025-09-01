package externalsorter

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log" //nolint:depguard // I am too lazy to use log/slog here. Sorry
	"os"
	"path/filepath"
	"sort"

	"github.com/S1riyS/wildberries-techschool/L2/10/internal/config"
	"github.com/S1riyS/wildberries-techschool/L2/10/internal/parser"
	"golang.org/x/sync/errgroup"
)

const (
	defaultChunkSize = 64 * 1024 * 1024 // 64 MB
	tempDirname      = "l2-10-sort-tmp"
)

// Line represents a line with its key for sorting.
type Line struct {
	Original string
	Key      string
}

// Sorter handles external sorting for files of any size.
type Sorter struct {
	config     *config.Config
	parser     *parser.Parser
	tempDir    string
	chunkSize  int
	chunkFiles []string
}

func MustNew(config *config.Config, parser *parser.Parser) *Sorter {
	tempDir := filepath.Join(os.TempDir(), tempDirname)
	err := os.MkdirAll(tempDir, 0600) // drw-------
	if err != nil {
		panic(err)
	}

	return &Sorter{
		config:    config,
		parser:    parser,
		tempDir:   tempDir,
		chunkSize: defaultChunkSize,
	}
}

// Sort performs external sort.
func (s *Sorter) Sort(input io.Reader, output io.Writer) error {
	// TODO: Implement sorting logic

	if err := s.createSortedChunks(input); err != nil {
		return err
	}

	if err := s.mergeSortedChunks(output); err != nil {
		return err
	}

	return nil
}

func (s *Sorter) createSortedChunks(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	var lines []Line
	currentSize := 0
	chunkCount := 0

	// Errorgroup to handle concurrent writing of chunks
	g, _ := errgroup.WithContext(context.Background())

	for scanner.Scan() {
		line := scanner.Text()
		parsed := s.parser.ParseLine(line)

		lines = append(lines, Line{Original: line, Key: parsed})
		currentSize += len(line)

		if currentSize >= s.chunkSize {
			// Write chunk concurrently
			g.Go(func() error {
				return s.writeChunk(lines, chunkCount)
			})

			// Reset
			lines = []Line{}
			currentSize = 0
			chunkCount++
		}
	}

	// Write last chunk
	if len(lines) > 0 {
		g.Go(func() error {
			return s.writeChunk(lines, chunkCount)
		})
	}

	// Wait for all chunks to be written
	if err := g.Wait(); err != nil {
		return err
	}

	return scanner.Err()
}

func (s *Sorter) writeChunk(lines []Line, chunkNum int) error {
	// Sort the chunk
	s.sortLines(lines)

	// Write to temporary file
	filename := filepath.Join(s.tempDir, fmt.Sprintf("chunk_%d.txt", chunkNum))
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err = writer.WriteString(line.Original + "\n")
		if err != nil {
			return err
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	s.chunkFiles = append(s.chunkFiles, filename)
	return nil
}

func (s *Sorter) sortLines(lines []Line) {
	sort.Slice(lines, func(i, j int) bool {
		compareResult := s.parser.Compare(lines[i].Key, lines[j].Key)
		if s.config.IsReverse {
			return compareResult > 0
		}
		return compareResult < 0
	})
}

func (s *Sorter) mergeSortedChunks(output io.Writer) error {
	if len(s.chunkFiles) == 0 {
		return nil
	}

	return s.multiWayMerge(output)
}

func (s *Sorter) multiWayMerge(output io.Writer) error {
	readers, files, err := s.openChunkFiles()
	if err != nil {
		return err
	}
	defer s.closeFiles(files)

	currentLines := s.initializeLines(readers)

	writer := bufio.NewWriter(output)
	defer writer.Flush()

	var lastLine string
	for {
		minIndex, minLine := s.findNextLine(currentLines)
		if minIndex == -1 {
			break // All chunks exhausted
		}

		if err = s.processLine(minLine, &lastLine, writer); err != nil {
			return err
		}

		s.advanceScanner(readers[minIndex], currentLines, minIndex)
	}

	return nil
}

func (s *Sorter) openChunkFiles() ([]*bufio.Scanner, []*os.File, error) {
	readers := make([]*bufio.Scanner, len(s.chunkFiles))
	files := make([]*os.File, len(s.chunkFiles))

	for i, filename := range s.chunkFiles {
		file, err := os.Open(filename)
		if err != nil {
			return nil, nil, err
		}
		files[i] = file
		readers[i] = bufio.NewScanner(file)
	}

	return readers, files, nil
}

func (s *Sorter) closeFiles(files []*os.File) {
	for _, file := range files {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}
}

// initializeLines reads the first line from each chunk file.
func (s *Sorter) initializeLines(readers []*bufio.Scanner) []*Line {
	currentLines := make([]*Line, len(readers))
	for i, scanner := range readers {
		if scanner.Scan() {
			line := scanner.Text()
			parsed := s.parser.ParseLine(line)
			currentLines[i] = &Line{Original: line, Key: parsed}
		}
	}
	return currentLines
}

// findNextLine returns the index of the line with the smallest key and the line itself.
func (s *Sorter) findNextLine(currentLines []*Line) (int, *Line) {
	minIndex := -1
	var minLine *Line

	for i, line := range currentLines {
		if line == nil {
			continue
		}

		if minLine == nil {
			minIndex = i
			minLine = line
			continue
		}

		compareResult := s.parser.Compare(line.Key, minLine.Key)
		if s.config.IsReverse {
			compareResult = -compareResult
		}

		if compareResult < 0 {
			minIndex = i
			minLine = line
		}
	}

	return minIndex, minLine
}

// processLine processes a line and writes it to the output.
func (s *Sorter) processLine(minLine *Line, lastLine *string, writer *bufio.Writer) error {
	// Handle same lines
	if minLine.Original == *lastLine && s.config.IsUnique {
		return nil
	}

	// Write line
	_, err := writer.WriteString(minLine.Original + "\n")
	if err != nil {
		return err
	}
	*lastLine = minLine.Original

	return nil
}

// advanceScanner advances the scanner for the given index.
func (s *Sorter) advanceScanner(scanner *bufio.Scanner, currentLines []*Line, index int) {
	if scanner.Scan() {
		line := scanner.Text()
		parsed := s.parser.ParseLine(line)
		currentLines[index] = &Line{Original: line, Key: parsed}
	} else {
		currentLines[index] = nil
	}
}

func (s *Sorter) Cleanup() {
	for _, chunkFile := range s.chunkFiles {
		err := os.Remove(chunkFile)
		if err != nil {
			log.Println(err)
		}
	}
}

// CheckIfSorted checks if the input is sorted.
// Returns true if sorted, false otherwise, along with the line number and contents of the unsorted line.
func (s *Sorter) CheckIfSorted(input io.Reader, output io.Writer) (bool, error) {
	scanner := bufio.NewScanner(input)

	lineNumber := 0
	var lastLine *Line
	for scanner.Scan() {
		lineNumber++

		line := scanner.Text()
		parsed := s.parser.ParseLine(line)
		currentLine := &Line{Original: line, Key: parsed}

		if lastLine == nil {
			lastLine = currentLine
			continue
		}

		comparisonResult := s.parser.Compare(currentLine.Key, lastLine.Key)
		if s.config.IsReverse {
			comparisonResult = -comparisonResult
		}

		// Check if sorted
		isSorted := true
		isSorted = isSorted && (comparisonResult != 0 || !s.config.IsUnique) // Uniqueness
		isSorted = isSorted && (comparisonResult > 0)                        // Ascending order

		if !isSorted {
			err := s.writeDisorderMessage(output, lineNumber, line)
			if err != nil {
				return false, err
			}
			return false, nil
		}

		lastLine = currentLine
	}

	return true, nil
}

func (s *Sorter) writeDisorderMessage(output io.Writer, lineNumber int, line string) error {
	// Determine filename
	var filenameToDisplay string
	if s.config.InputFile == "" {
		filenameToDisplay = "stdin"
	} else {
		filenameToDisplay = s.config.InputFile
	}

	// Write message
	message := fmt.Sprintf("%s:%d: disorder: %s\n", filenameToDisplay, lineNumber, line)
	_, err := output.Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}
