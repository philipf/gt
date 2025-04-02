package gtd

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
)

// addCmd represents the new command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new action to the inbox",
	Long:  `Adds a new action and adds it to the Kanban board. If a description is provided, it will be added to the action as a note.`,
	Run: func(_ *cobra.Command, _ []string) {
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
	gtdCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}

// Create a new action, this is an in memory representation of the action and will be persisted later
// Add a new todo to the kanban board
// Exist if no description was provided
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

	dueStr, err := getUserInput("Due date (optional, format: YYYY-MM-DD):", false)
	if err != nil {
		return err
	}

	// Parse the due date if it was provided
	var due time.Time
	if dueStr != "" {
		due, err = time.Parse("2006-01-02", dueStr)
		if err != nil {
			return err
		}
	}

	return addAction(title, description, &due)
}

func addAction(title string, description string, due *time.Time) error {
	action, err := gtd.CreateBasicAction(title, description, "cli")
	if err != nil {
		return err
	}

	containsDescription := strings.TrimSpace(action.Description) != ""

	err = gtd.AddToKanban(action.Title, containsDescription, due)
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
		return console.ReadLine()
	} else {
		lines, err := console.ReadMultiLine()
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

	llm, err := openai.New(
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

	// Based on the example: https://github.com/tmc/langchaingo/blob/main/examples/ernie-function-call-example/ernie_function_call_example.go

	ctx := context.Background()
	resp, err := llm.GenerateContent(ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman, prompt),
		},
		llms.WithTemperature(temperature),
		llms.WithFunctions(functions),
	)

	elapsedTime := time.Since(startTime)

	if err != nil {
		log.Fatal(err)
	}

	// In the updated API, the completion is a string that contains the function call response
	// We need to parse it directly
	var aiResponse map[string]string

	choice1 := resp.Choices[0]

	err = json.Unmarshal([]byte(choice1.FuncCall.Arguments), &aiResponse)

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
	shouldUse, err := console.ReadLine()
	if err != nil {
		return err
	}

	shouldUse = strings.TrimSpace(strings.ToLower(shouldUse))

	if shouldUse == "" || shouldUse == "y" {
		var duePtr *time.Time

		if aiResponse["dueDate"] != "" {
			due, err := time.Parse("2006-01-02", aiResponse["dueDate"])
			if err != nil {
				return err
			}
			duePtr = &due
		}

		return addAction(aiResponse["action"], aiResponse["summary"]+"\n\n## Original request\n"+input, duePtr)
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
					Description: "Due date of the action if it can be determined, the date format should be YYYY-MM-DD",
				},
			},
			Required: []string{"action", "summary"},
		},
	},
}
