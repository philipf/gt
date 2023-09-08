// Package files provides utility functions related to file handling.
package files

import (
	"strings"
)

const replacementChar = ' '

// ToValidFilename converts a title into a valid filename by replacing invalid characters.
//
// The function checks the title against a predefined set of invalid characters.
// Any character found in this set will be replaced by a predefined replacement character.
// The function also ensures that control characters (characters with ASCII values less than 32)
// are replaced to prevent any unforeseen issues.
// Finally, it trims any leading or trailing whitespace.
//
// Parameters:
//   - title: The input string that needs to be converted to a valid filename.
//
// Returns:
//   - string: A valid filename derived from the title.
func ToValidFilename(title string) string {
	// Define a set of characters that are invalid in filenames.
	invalidChars := map[rune]bool{
		'<': true, '>': true, ':': true, '"': true, '/': true,
		'\\': true, '|': true, '?': true, '*': true,
	}

	// Convert each character in the title:
	// If it's an invalid character, replace it with the replacement character.
	// If it's a control character (ASCII value < 32), replace it too.
	// Otherwise, let it be.
	validName := strings.Map(func(r rune) rune {
		if invalidChars[r] || r < 32 {
			return replacementChar
		}
		return r
	}, title)

	// Ensure that the filename does not have any leading or trailing whitespace.
	validName = strings.TrimSpace(validName)

	return validName
}
