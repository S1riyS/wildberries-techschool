package parser

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/S1riyS/wildberries-techschool/L2/10/internal/config"
)

type Parser struct {
	config *config.Config
}

func New(config *config.Config) *Parser {
	return &Parser{config: config}
}

// ParseLine parses a line and extracts the sort key.
func (p *Parser) ParseLine(line string) string {
	if p.config.IsIgnoreBlanks {
		line = strings.TrimSpace(line)
	}

	if p.config.Column > 0 {
		columns := strings.Split(line, "\t")
		if p.config.Column-1 < len(columns) {
			line = columns[p.config.Column-1]
		} else {
			line = ""
		}
	}

	return line
}

// Compare compares two lines based on the configuration.
func (p *Parser) Compare(a, b string) int {
	if p.config.IsNumeric {
		return p.compareNumeric(a, b)
	}
	if p.config.IsMonthSort {
		return p.compareMonth(a, b)
	}
	if p.config.IsHumanNumeric {
		return p.compareHuman(a, b)
	}
	return strings.Compare(a, b)
}

func (p *Parser) compareNumeric(a, b string) int {
	numA, errA := strconv.ParseFloat(a, 64)
	numB, errB := strconv.ParseFloat(b, 64)

	if errA != nil && errB != nil {
		return strings.Compare(a, b)
	}
	if errA != nil {
		return -1
	}
	if errB != nil {
		return 1
	}

	if numA < numB {
		return -1
	}
	if numA > numB {
		return 1
	}
	return 0
}

func (p *Parser) compareMonth(a, b string) int {
	monthA := p.parseMonth(a)
	monthB := p.parseMonth(b)

	// If both are 0, fall back to string comparison
	if monthA == 0 && monthB == 0 {
		return strings.Compare(a, b)
	}
	if monthA == 0 {
		return -1
	}
	if monthB == 0 {
		return 1
	}

	if monthA < monthB {
		return -1
	}
	if monthA > monthB {
		return 1
	}
	return 0
}

func (p *Parser) parseMonth(s string) time.Month {
	s = strings.ToLower(s)
	for prefix, month := range MonthNames {
		if strings.HasPrefix(s, prefix) {
			return month
		}
	}
	return 0
}

func (p *Parser) compareHuman(a, b string) int {
	valA := p.parseHumanSize(a)
	valB := p.parseHumanSize(b)

	if valA == -1 && valB == -1 {
		return strings.Compare(a, b)
	}
	if valA == -1 {
		return -1
	}
	if valB == -1 {
		return 1
	}

	if valA < valB {
		return -1
	}
	if valA > valB {
		return 1
	}
	return 0
}

func (p *Parser) parseHumanSize(s string) int64 {
	s = strings.TrimSpace(s)

	// Try to parse as plain number first
	if num, err := strconv.ParseInt(s, 10, 64); err == nil {
		return num
	}

	// Try to parse with suffix
	re := regexp.MustCompile(`^(\d+)([kKmMgGtT])?[bB]?$`)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return -1
	}

	num, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return -1
	}

	if len(matches) > 2 && matches[2] != "" {
		multiplier, exists := SizeSuffixes[matches[2]]
		if exists {
			num *= multiplier
		}
	}

	return num
}
