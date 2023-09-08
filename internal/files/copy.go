// Package files offers utilities for working with file operations.
package files

import (
	"io"
	"os"
)

// CopyFile copies the contents of the source file to a destination file.
//
// Parameters:
//   - src: The path to the source file.
//   - dst: The path to the destination file where the contents of the source file will be copied.
//     If a file with the same name exists, it will be overwritten.
//
// Returns:
//   - error: An error if any occurred during the file operations, otherwise nil.
func CopyFile(src, dst string) error {

	// Open the source file for reading.
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	// Ensure the source file is closed after operations are done.
	defer srcFile.Close()

	// Create or overwrite the destination file.
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	// Ensure the destination file is closed after operations are done.
	defer dstFile.Close()

	// Copy the contents from the source file to the destination file.
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Ensure that all changes are immediately written to the destination disk file.
	return dstFile.Sync()
}
