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
			_, err := output.Write([]byte(processedRow))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Cut) processRow(row string) (string, bool) {
	parts := strings.Split(row, c.cfg.Delimeter)
	if len(parts) == 0 {
		if c.cfg.IsOnlyDelimited {
			return "", false
		}
		return row, true
	}
	
	return "", true
}
