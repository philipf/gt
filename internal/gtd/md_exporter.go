package gtd

import (
	"errors"
	"os"
	"path/filepath"
	"text/template"

	"github.com/philipf/gt/internal/paths"
	"github.com/philipf/gt/internal/tasks"
)

const (
	InTemplatePath = "internal/gtd/templates/in.md"
)

func ActionToMd(action *Action, path string) error {
	return ActionsToMd(&[]Action{*action}, path)
}

func ActionsToMd(actions *[]Action, path string) error {
	tmpl, err := template.ParseFiles(InTemplatePath)
	if err != nil {
		return err
	}

	for _, action := range *actions {
		// check if the path exists and create it if not
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, 0755)
		}

		filename := filepath.Join(path, paths.TitleToFilename(action.Title)+".md")

		// If the file already exists return an error
		if _, err := os.Stat(filename); err == nil {
			return errors.New("file already exists: " + filename)
		}

		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		defer file.Close()

		err = tmpl.Execute(file, action)
		if err != nil {
			return err
		}
	}

	return nil
}

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
