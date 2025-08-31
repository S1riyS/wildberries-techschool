package unpacker

import (
	"errors"
	"strconv"
	"unicode"
)

// Unpack unpacks a string with repeating characters.
func Unpack(s string) (string, error) {
	var result []rune
	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		current := runes[i]

		// Escape character
		if isEscapeCharacter(current) {
			processed, newIndex, err := handleEscapeSequence(runes, i)
			if err != nil {
				return "", err
			}
			result = append(result, processed...)
			i = newIndex
			continue
		}

		// Digit
		if unicode.IsDigit(current) {
			processed, newIndex, err := handleDigitSequence(runes, i, result)
			if err != nil {
				return "", err
			}
			result = processed
			i = newIndex
			continue
		}

		// Regular character
		result = append(result, current)
	}

	return string(result), nil
}

func isEscapeCharacter(r rune) bool {
	return r == '\\'
}

// handleEscapeSequence handles the escape character sequence.
// It returns the escaped character and the index of the next character.
func handleEscapeSequence(runes []rune, currentIndex int) ([]rune, int, error) {
	if currentIndex+1 >= len(runes) {
		return nil, 0, errors.New("invalid escape sequence")
	}

	// Add the escaped character and skip it
	return []rune{runes[currentIndex+1]}, currentIndex + 1, nil
}

// handleDigitSequence handles the digit sequence.
// It returns the processed sequence and the index of the next character.
func handleDigitSequence(runes []rune, currentIndex int, currentResult []rune) ([]rune, int, error) {
	if currentIndex == 0 {
		return nil, 0, errors.New("invalid string: starts with a digit")
	}

	// Extract the full numeric sequence
	numStr, endIndex := extractNumberSequence(runes, currentIndex)
	count, err := strconv.Atoi(numStr)
	if err != nil {
		return nil, 0, err
	}

	result := applyRepetition(currentResult, count)
	return result, endIndex - 1, nil
}

// extractNumberSequence extracts the full numeric sequence from the string.
func extractNumberSequence(runes []rune, startIndex int) (string, int) {
	numStr := string(runes[startIndex])
	i := startIndex + 1

	for i < len(runes) && unicode.IsDigit(runes[i]) {
		numStr += string(runes[i])
		i++
	}

	return numStr, i
}

func applyRepetition(currentResult []rune, count int) []rune {
	if len(currentResult) == 0 {
		return currentResult
	}

	// Remove the last character
	if count == 0 {
		return currentResult[:len(currentResult)-1]
	}

	// Repeat the last character
	lastChar := currentResult[len(currentResult)-1]
	for k := 1; k < count; k++ {
		currentResult = append(currentResult, lastChar)
	}

	return currentResult
}
