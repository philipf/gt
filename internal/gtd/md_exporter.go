package gtd

import (
	"os"
	"text/template"

	"github.com/philipf/gt/internal/tasks"
)

func MapTasks(tasks []tasks.Task) ([]Action, error) {
	actions := make([]Action, len(tasks))
	for i, task := range tasks {
		actions[i] = Action{
			ID:           task.ID,
			ExternalID:   task.ExternalID,
			Title:        task.Title,
			Description:  task.Description,
			ExternalLink: task.ExternalLink,
			CreatedAt:    task.CreatedAt,
			ModifiedAt:   task.ModifiedAt,
		}
	}

	return actions, nil
}

func ExportToMd(actions []Action, path string) {

	tmpl, err := template.ParseFiles("internal/gtd/templates/in.md")

	for _, action := range actions {
		model := action
		if err != nil {
			panic(err)
		}

		tmpl.Execute(os.Stdout, model)
	}

	// Check if the path exists and create it if it doesn't

	// Create a file with the name of the action
	// Write the action to the file
	// Close the file

}
