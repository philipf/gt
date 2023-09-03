package gtd

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/philipf/gt/internal/tasks"
)

func MapTasks(tasks []tasks.Task) (*[]Action, error) {
	actions := make([]Action, len(tasks))
	for i, task := range tasks {
		actions[i] = Action{
			ID:           task.ID,
			ExternalID:   task.ExternalID,
			Title:        task.Title,
			Description:  task.Description,
			ExternalLink: task.ExternalLink,
			CreatedAt:    task.CreatedAt,
			UpdatedAt:    task.UpdatedAt,
		}
	}

	return &actions, nil
}

func ExportToMd(actions []Action, path string) {

	tmpl, err := template.ParseFiles("internal/gtd/templates/in.md")

	for _, action := range actions {
		model := action
		if err != nil {
			panic(err)
		}

		//tmpl.Execute(os.Stdout, model)

		// Create or open the file for writing

		// check if the path exists and create it if not
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, 0755)
		}

		// TODO: make sure the action.Title is a valid filename (no spaces, etc)
		filename := filepath.Join(path, action.Title+".md")
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		tmpl.Execute(file, model)
	}

}
