package grep

import (
	"testing"

	"github.com/S1riyS/wildberries-techschool/L2/10/internal/config"
)

func TestGrepFixedStringMatch(t *testing.T) {
	cfg := &config.Config{
		Pattern:    "test",
		IsFixedRow: true,
	}
	g := New(cfg, NewPrinter(cfg))

	// Test case sensitive match
	matches, start, end := g.fixedStringMatch("this is a test string")
	if !matches || start != 10 || end != 14 {
		t.Errorf("Expected match at position 10-14, got %t, %d, %d", matches, start, end)
	}

	// Test no match
	matches, _, _ = g.fixedStringMatch("no match here")
	if matches {
		t.Error("Expected no match")
	}
}

func TestGrepRegexpMatch(t *testing.T) {
	cfg := &config.Config{
		Pattern:    "t.st",
		IsFixedRow: false,
	}
	g := New(cfg, NewPrinter(cfg))

	// Test regex match
	matches, start, end := g.regexpMatch("this is a test string")
	if !matches || start != 10 || end != 14 {
		t.Errorf("Expected regex match at position 10-14, got %t, %d, %d", matches, start, end)
	}
}

func TestGrepCheckLineInverted(t *testing.T) {
	cfg := &config.Config{
		Pattern:    "test",
		IsInverted: true,
	}
	g := New(cfg, NewPrinter(cfg))

	// Test inverted match (should match when pattern is NOT found)
	matches, _, _ := g.checkLine(&Line{Text: "no match here"})
	if !matches {
		t.Error("Expected match for inverted pattern")
	}

	// Test inverted no match (should NOT match when pattern is found)
	matches, _, _ = g.checkLine(&Line{Text: "this is a test"})
	if matches {
		t.Error("Expected no match for inverted pattern when pattern exists")
	}
}
