package cut

import (
	"bufio"
	"io"
	"strings"

	"github.com/S1riyS/wildberries-techschool/L2/13/internal/config"
)

type Cut struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Cut {
	return &Cut{
		cfg: cfg,
	}
}

func (c *Cut) Run(input io.Reader, output io.Writer) error {
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		currentRow := scanner.Text()
		processedRow, shouldPrint := c.processRow(currentRow)
		if shouldPrint {
			_, err := output.Write([]byte(processedRow + "\n"))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Cut) processRow(row string) (string, bool) {
	parts := strings.Split(row, c.cfg.Delimiter)
	// Case: No delimiter
	if len(parts) == 1 {
		if c.cfg.IsOnlyDelimited {
			return "", false
		}
		return row, true
	}

	// Case: Has Delimiter
	var result []string
	// Iterate over all valid indexes
	left := max(c.cfg.Fields.MinIndex, 1)
	right := min(c.cfg.Fields.MaxIndex, len(parts))
	for i := left; i <= right; i++ {
		if c.cfg.Fields.IsInRange(i) {
			index := i - 1 // ! switch to 0-index (config is 1-index)
			result = append(result, parts[index])
		}
	}

	if len(result) > 0 {
		return strings.Join(result, c.cfg.Delimiter), true
	}
	return "", true
}
