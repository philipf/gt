package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/philipf/gt/internal/gtd"
	"github.com/philipf/gt/internal/settings"
	"github.com/spf13/cobra"
)

// gtdCmd represents the action command
var gtdCmd = &cobra.Command{
	Use:   "gtd",
	Short: "Create a new GTD action",
	Long: `Create a new GTD action. This will create a new action in the inbox and add a new todo to the kanban board.
	If now description is provided, only the todo will be added to the kanban board.`,

	Run: func(cmd *cobra.Command, args []string) {
		err := promptForAction()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gtdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// actionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// actionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func promptForAction() error {
	// User input
	title, err := getUserInput("Action title:", false)
	if err != nil {
		return err
	}
	if title == "" {
		return errors.New("title cannot be empty")
	}

	description, err := getUserInput("Description (optional):", true)
	if err != nil {
		return err
	}

	// Create a new action, this is an in memory representation of the action and will be persisted later
	action, err := gtd.CreateBasicAction(title, description, "cli")
	if err != nil {
		return err
	}

	containsDescription := strings.TrimSpace(action.Description) != ""

	// Add a new todo to the kanban board
	err = addToKanban(action.Title, containsDescription)
	if err == nil {
		fmt.Println("To-do added to kanban board")
	} else {
		return err
	}

	// Exist if no description was provided
	if !containsDescription {
		return nil
	}

	err = addDescriptionNote(action)
	if err == nil {
		fmt.Println("To-do and description added")
	} else {
		return err
	}

	return nil
}

func addDescriptionNote(action *gtd.Action) error {
	inTemplate := settings.GetInTemplatePath()
	if _, err := os.Stat(inTemplate); os.IsNotExist(err) {
		return err
	}

	inboxPath := settings.GetKanbanInboxPath()

	err := gtd.ActionToMd(action, inTemplate, inboxPath)
	if err != nil {
		return err
	}
	return nil
}

func addToKanban(todo string, withLink bool) error {
	path := settings.GetKanbanBoardPath()
	err := gtd.InsertTodo(path, settings.GetKanbanInColumn(), todo, withLink)
	return err
}

func getUserInput(prompt string, allowMultiLine bool) (string, error) {
	fmt.Println(prompt)

	if !allowMultiLine {
		return readSingleLineInput()
	} else {
		lines, err := readMultiLineInput()
		if err != nil {
			return "", err
		}
		result := strings.Join(lines, "\n")
		return result, nil
	}
}

func readSingleLineInput() (string, error) {
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
func readMultiLineInput() ([]string, error) {

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
