/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/philipf/gt/internal/console"
	"github.com/philipf/gt/internal/gtd"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if UseAi {
			//err := promptForActionUsingAi()
			fmt.Println("AI not implemented yet")
			err = nil
		} else {
			err = promptForAction()
		}
		cobra.CheckErr(err)
	},
}

func init() {
	gtdCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func promptForAction() error {
	title, err := getUserInput("Action title:", false)
	if err != nil {
		return err
	}

	if title == "" {
		return errors.New("title cannot be empty")
	}

	description, err := getUserInput("Description (optional, end capture with a full stop on a new line):", true)
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
	err = gtd.AddToKanban(action.Title, containsDescription)
	if err == nil {
		fmt.Println("To-do added to kanban board")
	} else {
		return err
	}

	// Exist if no description was provided
	if !containsDescription {
		return nil
	}

	err = gtd.AddDescriptionNote(action)
	if err == nil {
		fmt.Println("To-do and description added")
	} else {
		return err
	}

	return nil
}

func getUserInput(prompt string, allowMultiLine bool) (string, error) {
	fmt.Println(prompt)

	if !allowMultiLine {
		return console.ReadSingleLineInput()
	} else {
		lines, err := console.ReadMultiLineInput()
		if err != nil {
			return "", err
		}
		result := strings.Join(lines, "\n")
		return result, nil
	}
}
