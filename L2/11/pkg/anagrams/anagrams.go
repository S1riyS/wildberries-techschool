package anagrams

import (
	"slices"
	"sort"
	"strings"
)

// Find returns a map of anagram groups
//
// Time complexity: O(n*m*log(m)), where m is avgerage word length.
func Find(data []string) map[string][]string {
	// Convert data to lowercase
	loweredData := make([]string, len(data))
	for i, word := range data {
		loweredData[i] = strings.ToLower(word)
	}

	firstOccurrence := make(map[string]string)
	anagramGroups := make(map[string][]string)

	for _, word := range loweredData {
		// Sort runes of current word
		runes := []rune(word)
		slices.Sort(runes)
		sortedWord := string(runes)

		// Store the first occurrence of each anagram pattern
		_, ok := anagramGroups[sortedWord]
		if !ok {
			firstOccurrence[sortedWord] = word
		}

		// Add word to anagram group
		anagramGroups[sortedWord] = append(anagramGroups[sortedWord], word)
	}

	// Assemble resulting map
	result := make(map[string][]string)
	for sortedKey, group := range anagramGroups {
		// If there were found any anagrams, add to result
		if len(group) > 1 {
			initialString := firstOccurrence[sortedKey]
			sort.Strings(group)
			result[initialString] = group
		}
	}

	return result
}
