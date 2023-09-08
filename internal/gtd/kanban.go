// Package gtd (Getting Things Done) offers utilities to manage a Kanban-style TODO list using the Obsidian Kanban plugin.
package gtd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	// In the Kanban board, headings are prefixed with ## and a space.
	headingDiscriminator = "## "

	// To-do items are prefixed with either a - [ ] or - [x] depending on their completion status.
	todoDiscriminator     = "- [ ] "
	todoDiscriminatorDone = "- [x] "
)

// findLastToDo locates the last to-do item under a specified heading in a list of lines.
// Returns the index of the last to-do or the heading if no to-dos are found under it.
func findLastToDo(lines []string, heading string) int {
	headingToFind := headingDiscriminator + heading

	insideHeading := false
	headingIdx := -1
	lastTodo := -1

	// A state machine to find the last to-do in the specified heading.
	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Check if we're inside the desired heading section.
		if insideHeading {
			// If we find a to-do item, record its index.
			if strings.HasPrefix(line, todoDiscriminator) || strings.HasPrefix(line, todoDiscriminatorDone) {
				lastTodo = i
				continue
			} else if strings.HasPrefix(line, headingDiscriminator) {
				// If we encounter another heading, the previous heading section has ended.
				// Return the index of the last found to-do or the heading itself.
				return max(lastTodo, headingIdx)
			}
		} else if line == headingToFind {
			// Mark that we've entered the desired heading section.
			insideHeading = true
			headingIdx = i
			lastTodo = headingIdx
			continue
		}
	}

	return lastTodo
}

// readFile reads the content of a file at a given path and returns its lines.
func readFile(path string) ([]string, error) {
	// Open the file for reading.
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	// Scan each line and append it to the lines slice.
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Return the collected lines and any scanning error.
	return lines, scanner.Err()
}

// writeFile writes a slice of lines to a file at the specified path.
func writeFile(path string, lines []string) error {
	// Join lines with newline and write to the file.
	err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
	return err
}

// InsertTodo inserts a new to-do item under a specified heading in a file.
// If withLink is true, the to-do item is formatted as a link.
func InsertTodo(path, heading, todo string, withLink bool) error {
	if withLink {
		todo = fmt.Sprintf("[[%s]]", todo)
	}

	lines, err := readFile(path)
	if err != nil {
		return err
	}

	// Find the position to insert the new to-do.
	lastTodo := findLastToDo(lines, heading)
	if lastTodo == -1 {
		return fmt.Errorf("could not find heading [%s] in the file [%s]", heading, path)
	}

	// Insert the new to-do at the correct position.
	todo = todoDiscriminator + todo
	lines = append(lines[:lastTodo+1], append([]string{todo}, lines[lastTodo+1:]...)...)

	// Write the updated lines back to the file.
	err = writeFile(path, lines)
	return err
}

// max returns the maximum of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
