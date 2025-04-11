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

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
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
}

// Create a new action, this is an in memory representation of the action and will be persisted later
// Add a new todo to the kanban board
// Exist if no description was provided
func promptForAction() error {
	// Define variables to store form values
	var title string
	var description string
	var dueStr string

	// Default dueStr to today's date if empty
	if strings.TrimSpace(dueStr) == "" {
		dueStr = time.Now().Format("2006-01-02")
	}

	// Create the form with Huh library
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Action title").
				Value(&title).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("title cannot be empty")
					}
					return nil
				}),

			huh.NewText().
				Title("Description (optional)").
				CharLimit(0).
				Value(&description),

			huh.NewInput().
				Title("Due date (optional, format: YYYY-MM-DD)").
				Value(&dueStr).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return nil // Empty is valid
					}
					_, err := time.Parse("2006-01-02", s)
					if err != nil {
						return errors.New("due date must be in YYYY-MM-DD format")
					}
					return nil
				}),
		),
	)

	// Run the form
	err := form.Run()
	if err != nil {
		return err
	}

	// Parse the due date if it was provided
	var due time.Time
	if dueStr != "" {
		due, err = time.Parse("2006-01-02", dueStr)
		if err != nil {
			return err // Should not happen due to validation
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

func promptForActionUsingAi() error {
	var input string

	// Create a form for AI input
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Enter input to create an action from:").
				Value(&input),
		),
	).WithKeyMap(createDefaultKeyMap())

	err := form.Run()
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

	data := map[string]any{
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

	// Display processing time
	fmt.Printf("Operation completed in %.2f seconds\n", elapsedTime.Seconds())

	// Create editable form with the AI-generated values
	title := aiResponse["action"]
	summary := aiResponse["summary"]
	dueDate := aiResponse["dueDate"]
	// Default dueDate to today's date if empty
	if strings.TrimSpace(dueDate) == "" {
		dueDate = time.Now().Format("2006-01-02")
	}
	confirmAdd := true

	// Create a confirmation form with the AI response, allowing edits
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("AI Generated Action").
				Description("Review and edit if needed:"),

			huh.NewInput().
				Title("Action Title").
				Value(&title).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("title cannot be empty")
					}
					return nil
				}),

			huh.NewText().
				Title("Description").
				Value(&summary).
				CharLimit(0),

			huh.NewInput().
				Title("Due Date (YYYY-MM-DD)").
				Value(&dueDate).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return nil // Empty is valid
					}
					_, err := time.Parse("2006-01-02", s)
					if err != nil {
						return errors.New("due date must be in YYYY-MM-DD format")
					}
					return nil
				}),

			huh.NewConfirm().
				Title("Add this action?").
				Value(&confirmAdd),
		),
	)

	err = form.Run()
	if err != nil {
		return err
	}

	if confirmAdd {
		var duePtr *time.Time

		if dueDate != "" {
			due, err := time.Parse("2006-01-02", dueDate)
			if err != nil {
				return err // Should not happen due to validation
			}
			duePtr = &due
		}

		// Include the original input in the description
		fullDescription := summary
		if !strings.Contains(summary, "Original request") {
			fullDescription += "\n\n## Original request\n" + input
		}

		return addAction(title, fullDescription, duePtr)
	} else {
		return promptForAction()
	}
}

func createDefaultKeyMap() *huh.KeyMap {
	keyMap := huh.NewDefaultKeyMap()
	// Change NewLine to use Enter key
	keyMap.Text.NewLine = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "new line"),
	)
	// Change Submit to use Control+Enter
	keyMap.Text.Submit = key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "submit"),
	)
	return keyMap
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
