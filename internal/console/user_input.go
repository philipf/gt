package console

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Prompt(prompt string) (string, error) {
	fmt.Print(prompt)
	return ReadLine()
}

func PromptInt64(prompt string) (int64, error) {
	fmt.Print(prompt)
	l, err := ReadLine()
	if err != nil {
		return 0, err
	}

	// convert string to int64
	var i int64
	i, err = strconv.ParseInt(l, 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func PromptMultiLine(prompt string) ([]string, error) {
	fmt.Print(prompt)
	return ReadMultiLine()
}

func ReadLine() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		return input, nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

// Multi line logic allows the user to enter multiple lines of text and ends when the users enters a full stop on a new line
func ReadMultiLine() ([]string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		// If is the first line and the user didn't enter anything, return an empty string, this a quick way to exit the prompt
		if len(lines) == 0 && input == "" {
			return lines, nil
		}

		if input == "." {
			break
		}

		lines = append(lines, input)
	}

	if err := scanner.Err(); err != nil {
		return lines, err
	}

	return lines, nil
}
