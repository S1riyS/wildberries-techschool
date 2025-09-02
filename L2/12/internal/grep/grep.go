package grep

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/S1riyS/wildberries-techschool/L2/10/internal/config"
	"github.com/S1riyS/wildberries-techschool/L2/10/internal/datastructure"
)

type Grep struct {
	config           *config.Config
	printer          *Printer
	beforeContext    *datastructure.SlidingWindow[*Line] // Slinding window of lines before the current line
	afterContextSize int                                 // The size of the context window after the current line
}

func New(config *config.Config, printer *Printer) *Grep {
	beforeContextSize := max(config.RowsBefore, config.RowsAround)
	afterContextSize := max(config.RowsAfter, config.RowsAround)

	return &Grep{
		config:           config,
		printer:          printer,
		beforeContext:    datastructure.NewSlidingWindow[*Line](beforeContextSize),
		afterContextSize: afterContextSize,
	}
}

// Run runs the grep command. Returns the number of matches.
func (g *Grep) Run(input io.Reader, output io.Writer) int {
	scanner := bufio.NewScanner(input)

	// Initial values
	counter := 0
	rowsToPrintAfter := 0
	rowNumber := 0

	for scanner.Scan() {
		rowNumber++
		currentLine := NewLine(rowNumber, scanner.Text())

		matches, start, end := g.checkLine(currentLine)

		// Handle pattern match
		if matches {
			counter++
			rowsToPrintAfter = g.afterContextSize
			currentLine.SetMatch(start, end)
		}

		// Print lines
		if !g.config.IsOnlyRowsCount {
			if matches {
				// Previous lines
				for previousLine := range g.beforeContext.Iterator() {
					g.printer.Print(previousLine, output)
				}
				// Current line
				g.printer.Print(currentLine, output)
			} else if rowsToPrintAfter > 0 {
				// If line doesn't match but should be printed due to -A or -C option
				g.printer.Print(currentLine, output)
				rowsToPrintAfter--
				if rowsToPrintAfter == 0 {
					// Print group separator
					g.printer.PrintGroupSeparator(output)
				}
			}

			g.beforeContext.Add(currentLine) // Add line to sliding window
		}
	}

	// Print count of matching lines
	if g.config.IsOnlyRowsCount {
		_, _ = fmt.Fprintf(output, "%d\n", counter)
	}

	return counter
}

func (g *Grep) checkLine(line *Line) (bool, int, int) {
	var matches bool
	var start, end int

	switch {
	case g.config.IsFixedRow:
		matches, start, end = g.fixedStringMatch(line.Text)
	default:
		matches, start, end = g.regexpMatch(line.Text)
	}

	// Invert result if needed
	if g.config.IsInverted {
		matches = !matches
		// For reversed matches, start and end indexes doesn't make sense
		if matches {
			start, end = -1, -1
		}
	}

	return matches, start, end
}

func (g *Grep) fixedStringMatch(value string) (bool, int, int) {
	var pattern, searchValue string

	if g.config.IsIgnoreRegister {
		pattern = strings.ToLower(g.config.Pattern)
		searchValue = strings.ToLower(value)
	} else {
		pattern = g.config.Pattern
		searchValue = value
	}

	index := strings.Index(searchValue, pattern)
	if index == -1 {
		return false, -1, -1
	}

	return true, index, index + len(pattern)
}

func (g *Grep) regexpMatch(value string) (bool, int, int) {
	pattern := g.config.Pattern
	if g.config.IsIgnoreRegister {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		// Fallback to simple string matching on regexp compile error
		return g.fixedStringMatch(value)
	}

	// Find the first match
	loc := re.FindStringIndex(value)
	if loc == nil {
		return false, -1, -1
	}

	return true, loc[0], loc[1]
}
