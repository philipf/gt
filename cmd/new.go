/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/philipf/gt/internal/console"
	"github.com/philipf/gt/internal/gtd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tmc/langchaingo/jsonschema"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
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
			err = promptForActionUsingAi()
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
	// Add a new todo to the kanban board
	// Exist if no description was provided
	return addAction(title, description)
}

func addAction(title string, description string) error {
	action, err := gtd.CreateBasicAction(title, description, "cli")
	if err != nil {
		return err
	}

	containsDescription := strings.TrimSpace(action.Description) != ""

	err = gtd.AddToKanban(action.Title, containsDescription)
	if err == nil {
		fmt.Println("To-do added to kanban board")
	} else {
		return err
	}

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

func promptForActionUsingAi() error {
	input, err := getUserInput("Enter input to create an action from:", true)

	if err != nil {
		return err
	}

	fmt.Printf("Processing using [%s]....\n", viper.GetString("ai.openAiModel"))

	llm, err := openai.NewChat(
		openai.WithModel(viper.GetString("ai.openAiModel")),
		openai.WithToken(viper.GetString("ai.openAiKey")),
	)

	if err != nil {
		log.Fatal(err)
	}

	templateStr := viper.GetString("gtd.actionPrompt.user")
	temperature := viper.GetFloat64("gtd.actionPrompt.temperature")

	t := template.Must(template.New("gtdTemplate").Parse(templateStr))

	data := map[string]interface{}{
		"name":        viper.GetString("personal.name"),
		"surname":     viper.GetString("personal.surname"),
		"company":     viper.GetString("personal.company"),
		"currentDate": time.Now().Format(time.RFC3339),
		"input":       input,
	}

	var resultBuffer bytes.Buffer

	err = t.Execute(&resultBuffer, data)
	if err != nil {
		log.Fatal(err)
	}

	prompt := resultBuffer.String()

	startTime := time.Now()

	ctx := context.Background()
	completion, err := llm.Call(ctx, []schema.ChatMessage{
		schema.HumanChatMessage{Content: prompt}, // For some reason specifying a system message causes GPT-4 to ignore the FunctionCall??
	}, llms.WithTemperature(temperature),
		llms.WithFunctions(functions),
	)

	elapsedTime := time.Since(startTime)

	if err != nil {
		log.Fatal(err)
	}

	if completion.FunctionCall == nil {
		log.Fatal("No function call returned")
	}

	// Store the the function call arguments in a map by first parsing it to json
	var aiResponse map[string]string

	err = json.Unmarshal([]byte(completion.FunctionCall.Arguments), &aiResponse)

	if err != nil {
		log.Fatal(err)
	}

	// Print the action and summary to the console
	fmt.Println("--------------------------------------------")
	fmt.Println("Action  : ", aiResponse["action"])
	fmt.Println("Summary : ", aiResponse["summary"])
	fmt.Println("Due date: ", aiResponse["dueDate"])
	fmt.Println("--------------------------------------------")
	fmt.Printf("Operation took %.2f seconds\n", elapsedTime.Seconds())

	// Ask the user if they want to add the action
	fmt.Println("Do you want to add this action? [y]/n")
	shouldUse, err := console.ReadSingleLineInput()
	if err != nil {
		return err
	}

	shouldUse = strings.TrimSpace(strings.ToLower(shouldUse))

	if shouldUse == "" || shouldUse == "y" {
		return addAction(aiResponse["action"], aiResponse["summary"]+"\n\n## Original request\n"+input)
	} else {
		return promptForAction()
	}
}

var functions = []llms.FunctionDefinition{
	{
		Name:        "getAction",
		Description: "Get an action, summary and due date from the user",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"action": {
					Type:        jsonschema.String,
					Description: "Action to be taken, in less than 10 words",
				},
				"summary": {
					Type:        jsonschema.String,
					Description: "A summary of the user input, proof read and edited to be concise to less than 200 words",
				},
				"dueDate": {
					Type:        jsonschema.String,
					Description: "Due date of the action if it can be determined",
				},
			},
			Required: []string{"action", "summary"},
		},
	},
}
