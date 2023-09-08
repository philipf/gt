package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/philipf/gt/internal/gtd"
	"github.com/spf13/cobra"
)

// gtdCmd represents the action command
var gtdCmd = &cobra.Command{
	Use:   "gtd",
	Short: "Create a new GTD action",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		dt := time.Now()
		externalId := ""
		externalLink := ""

		title, err := getInputWithPrompt("Please enter a title:")
		if err != nil {
			panic(err)
		}

		description, err := getInputWithPrompt("Please enter a description:")
		if err != nil {
			panic(err)
		}

		action, err := gtd.CreateAction(externalId, title, description, externalLink, dt, dt, gtd.In)
		if err != nil {
			panic(err)
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		inTemplate := filepath.Join(homeDir, "gt", "in.md")

		// check if inTemplate exists
		if _, err := os.Stat(inTemplate); os.IsNotExist(err) {
			panic(err)
		}

		// Add a new todo to the kanban board
		kanbanPath := "G:/My Drive/SecondBrain/_GTD/_Board.md"
		todo := fmt.Sprintf("[[%s]]", action.Title)
		gtd.InsertTodo(kanbanPath, "In", todo)

		// Add a new action, in the form a markdown file, to the inbox
		inboxPath := "G:/My Drive/SecondBrain/_GTD/Inbox"

		err = gtd.ActionToMd(action, inTemplate, inboxPath)
		if err != nil {
			panic(err)
		}
	},
}

func getInputWithPrompt(prompt string) (string, error) {
	fmt.Println(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()
		return input, nil
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
		return "", err
	}
	return "", nil
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
