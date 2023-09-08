// Package files provides utilities for working with files.
package files

import (
	"bytes"
	"os"
)

// AreFilesEqual checks if the contents of two files are identical.
//
// Parameters:
//   - file1: The path to the first file.
//   - file2: The path to the second file.
//
// Returns:
//   - bool: True if the contents of both files are identical, otherwise false.
//   - error: An error if any occurred while reading the files, otherwise nil.
func AreFilesEqual(file1, file2 string) (bool, error) {

	// Read the contents of the first file.
	content1, err1 := os.ReadFile(file1)
	if err1 != nil {
		return false, err1
	}

	// Read the contents of the second file.
	content2, err2 := os.ReadFile(file2)
	if err2 != nil {
		return false, err2
	}

	// Compare the contents of the two files and return the result.
	return bytes.Equal(content1, content2), nil
}
