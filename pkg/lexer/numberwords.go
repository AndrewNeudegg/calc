package lexer

import (
	"strings"
)

// NumberWord represents a word that can be converted to a number
type NumberWord struct {
	Word  string
	Value float64
}

// GetNumberWords returns the number words for a given locale
func GetNumberWords(locale string) map[string]float64 {
	// Default to en_GB, but works for en_US too
	if strings.HasPrefix(locale, "en_") || locale == "" {
		return enNumberWords
	}
	return enNumberWords
}

// English number words (en_GB and en_US)
var enNumberWords = map[string]float64{
	// Basic numbers 0-20
	"zero":       0,
	"one":        1,
	"two":        2,
	"three":      3,
	"four":       4,
	"five":       5,
	"six":        6,
	"seven":      7,
	"eight":      8,
	"nine":       9,
	"ten":        10,
	"eleven":     11,
	"twelve":     12,
	"thirteen":   13,
	"fourteen":   14,
	"fifteen":    15,
	"sixteen":    16,
	"seventeen":  17,
	"eighteen":   18,
	"nineteen":   19,
	"twenty":     20,
	
	// Tens
	"thirty":     30,
	"forty":      40,
	"fifty":      50,
	"sixty":      60,
	"seventy":    70,
	"eighty":     80,
	"ninety":     90,
	
	// Scale words
	"hundred":    100,
	"thousand":   1000,
	"million":    1000000,
	"billion":    1000000000,
	"trillion":   1000000000000,
	
	// Fractions (for consistency)
	"half":       0.5,
	"quarter":    0.25,
}

// Connector words that should be ignored
var connectorWords = map[string]bool{
	"and": true,
	"a":   true,
	"an":  true,
}

// ParseNumberWords attempts to parse a sequence of words as a number
// Returns the number value and true if successful, 0 and false otherwise
func ParseNumberWords(words []string, locale string) (float64, bool) {
	if len(words) == 0 {
		return 0, false
	}
	
	numberWords := GetNumberWords(locale)
	
	// Single word case
	if len(words) == 1 {
		word := strings.ToLower(words[0])
		if val, ok := numberWords[word]; ok {
			return val, true
		}
		return 0, false
	}
	
	// Multi-word number parsing
	var total float64
	var current float64
	
	for i := 0; i < len(words); i++ {
		word := strings.ToLower(words[i])
		
		// Skip connector words
		if connectorWords[word] {
			continue
		}
		
		val, exists := numberWords[word]
		if !exists {
			return 0, false // Not a valid number word sequence
		}
		
		// Handle scale words (hundred, thousand, million, etc.)
		if val >= 100 {
			if current == 0 {
				current = 1 // "hundred" means "one hundred"
			}
			current *= val
			
			// If this is thousand/million/billion, add to total
			if val >= 1000 {
				total += current
				current = 0
			}
		} else {
			current += val
		}
	}
	
	total += current
	
	if total > 0 {
		return total, true
	}
	
	return 0, false
}

// IsNumberWord checks if a single word is a number word
func IsNumberWord(word string, locale string) bool {
	numberWords := GetNumberWords(locale)
	_, exists := numberWords[strings.ToLower(word)]
	if exists {
		return true
	}
	return connectorWords[strings.ToLower(word)]
}
