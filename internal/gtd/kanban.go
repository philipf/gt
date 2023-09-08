package gtd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	// In the Kanban board, headings are prefixed with ## and a space
	headingDiscriminator = "## "

	// To do items are prefixed with either a - [ ] or - [x] depending on whether they are done or not
	todoDiscriminator     = "- [ ] "
	todoDiscriminatorDone = "- [x] "
)

func findLastToDo(lines []string, heading string) int {
	headingToFind := headingDiscriminator + heading

	insideHeading := false
	headingIdx := -1
	lastTodo := -1

	// Basic state machine to find the last to do in the In heading
	for i, line := range lines {
		line = strings.TrimSpace(line)

		if insideHeading {
			// The current state is inside the In heading, so we need to check if we are still inside it
			// and record the last todo, if any

			if strings.HasPrefix(line, todoDiscriminator) || strings.HasPrefix(line, todoDiscriminatorDone) {
				// A to do was found, but we don't know if it's the last one yet
				lastTodo = i
				continue
			} else if strings.HasPrefix(line, headingDiscriminator) {
				// Found a new heading, so we are no longer inside the heading we were looking for and we can stop
				// Return the last todo, if any was found otherwise the last line of the heading we were looking for
				if lastTodo == -1 {
					return headingIdx
				} else {
					return lastTodo
				}
			}
		} else if line == headingToFind {
			// We found the heading we were looking for, so we set the state to inside the heading
			insideHeading = true
			headingIdx = i
			lastTodo = headingIdx // this is a safety in case only the In heading is found
			continue
		}
	}

	return lastTodo
}

func readFile(path string) ([]string, error) {
	// Open the file for reading
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// Make sure we close the file when we're done
	defer file.Close()

	// Create a new scanner that reads from the file
	scanner := bufio.NewScanner(file)

	// Create a slice of strings to store the lines in
	var lines []string

	// Use the scanner to loop through all the lines in the file
	for scanner.Scan() {
		// Add the line to the lines slice
		lines = append(lines, scanner.Text())
	}

	// Return the lines and any error that happened
	return lines, scanner.Err()
}

func writeFile(path string, lines []string) error {
	err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return err
	}

	return nil
}

func InsertTodo(path, heading, todo string, withLink bool) error {
	if withLink {
		todo = fmt.Sprintf("[[%s]]", todo)
	}

	lines, err := readFile(path)
	if err != nil {
		return err
	}

	lastTodo := findLastToDo(lines, heading)
	if lastTodo == -1 {
		return fmt.Errorf("could not find heading [%s] in the file [%s]", heading, path)
	}

	// Add the new to do in the correct location in the slice
	todo = todoDiscriminator + todo
	lines = append(lines[:lastTodo+1], append([]string{todo}, lines[lastTodo+1:]...)...)

	err = writeFile(path, lines)
	if err != nil {
		return err
	}

	return nil
}
