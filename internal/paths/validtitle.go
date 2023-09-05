package paths

import (
	"strings"
)

const replacementChar = ' '

func TitleToFilename(title string) string {
	// Define a set of invalid characters for filenames.
	invalidChars := map[rune]bool{
		'<': true, '>': true, ':': true, '"': true, '/': true,
		'\\': true, '|': true, '?': true, '*': true,
	}

	// Replace each invalid character with an underscore.
	cleanTitle := strings.Map(func(r rune) rune {
		if invalidChars[r] || r < 32 {
			return replacementChar
		}
		return r
	}, title)

	// Remove any leading or trailing whitespace.
	cleanTitle = strings.TrimSpace(cleanTitle)

	// Replace spaces with underscores for better filename readability.
	//cleanTitle = strings.ReplaceAll(cleanTitle, " ", string(replacementChar))

	return cleanTitle
}
