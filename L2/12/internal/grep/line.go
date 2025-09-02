package grep

type Line struct {
	IsPrinted  bool // Whether the line has been printed
	IsMatched  bool // Whether the line matches the pattern
	Number     int
	Text       string
	MatchStart int // Start index of the match (-1 if no match)
	MatchEnd   int // End index of the match (-1 if no match)
}

func NewLine(number int, value string) *Line {
	return &Line{
		IsPrinted:  false,
		IsMatched:  false,
		Number:     number,
		Text:       value,
		MatchStart: -1,
		MatchEnd:   -1,
	}
}

func (line *Line) SetMatch(start, end int) {
	line.IsMatched = true
	line.MatchStart = start
	line.MatchEnd = end
}
