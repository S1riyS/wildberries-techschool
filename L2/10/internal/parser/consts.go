package parser

import "time"

//nolint:gochecknoglobals // will not be changed
var MonthNames = map[string]time.Month{
	"jan": time.January, "feb": time.February, "mar": time.March,
	"apr": time.April, "may": time.May, "jun": time.June,
	"jul": time.July, "aug": time.August, "sep": time.September,
	"oct": time.October, "nov": time.November, "dec": time.December,
}

//nolint:gochecknoglobals,mnd // intentional bit shifts for size calculation
var SizeSuffixes = map[string]int64{
	"k": 1 << 10, "m": 1 << 20, "g": 1 << 30, "t": 1 << 40,
	"K": 1 << 10, "M": 1 << 20, "G": 1 << 30, "T": 1 << 40,
}
