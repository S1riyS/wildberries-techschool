package grep

import (
	"fmt"
	"io"

	"github.com/S1riyS/wildberries-techschool/L2/10/internal/config"
)

type Color string

const (
	Reset  Color = "\033[0m"
	Red    Color = "\033[31m"
	Green  Color = "\033[32m"
	Yellow Color = "\033[33m"
	Blue   Color = "\033[34m"
	Purple Color = "\033[35m"
	Cyan   Color = "\033[36m"
	Gray   Color = "\033[37m"
	White  Color = "\033[97m"
)

const (
	defaultGroupSeparator = "--"
)

type Printer struct {
	config *config.Config
}

func Colorize(text string, colorCode Color) string {
	return string(colorCode) + text + string(Reset)
}

func NewPrinter(config *config.Config) *Printer {
	return &Printer{
		config: config,
	}
}

// Print prints a line to the output.
// If the line has already been printed, it does nothing.
func (p *Printer) Print(line *Line, output io.Writer) {
	if line.IsPrinted {
		return
	}

	// Print row number
	if p.config.IsPrintRowNumber {
		_, _ = fmt.Fprintf(output, Colorize("%d", Green), line.Number)

		// Print delimiter
		var delimiter string
		if line.IsMatched {
			delimiter = ":"
		} else {
			delimiter = "-"
		}
		_, _ = fmt.Fprint(output, Colorize(delimiter, Cyan))
	}

	// Print line
	if line.IsMatched && line.MatchStart != -1 && line.MatchEnd != -1 {
		// Colorize match
		beforeMatch := line.Text[:line.MatchStart]
		matchText := line.Text[line.MatchStart:line.MatchEnd]
		afterMatch := line.Text[line.MatchEnd:]
		coloredLine := beforeMatch + Colorize(matchText, Red) + afterMatch
		_, _ = fmt.Fprint(output, coloredLine)
	} else {
		_, _ = fmt.Fprint(output, line.Text)
	}

	// Print newline
	_, _ = fmt.Fprint(output, "\n")

	line.IsPrinted = true
}

func (p *Printer) PrintGroupSeparator(output io.Writer) {
	_, _ = fmt.Fprint(output, Colorize(defaultGroupSeparator, Cyan))
	_, _ = fmt.Fprint(output, "\n")
}
