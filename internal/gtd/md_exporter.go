package gtd

import (
	"errors"
	"os"
	"path/filepath"
	"text/template"

	"github.com/philipf/gt/internal/paths"
	"github.com/philipf/gt/internal/tasks"
)

func ActionToMd(action *Action, templateFile, outputPath string) error {
	return ActionsToMd(&[]Action{*action}, templateFile, outputPath)
}

func ActionsToMd(actions *[]Action, templateFile, outputPath string) error {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	for _, action := range *actions {
		// check if the path exists and create it if not
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			os.Mkdir(outputPath, 0755)
		}

		filename := filepath.Join(outputPath, paths.TitleToFilename(action.Title)+".md")

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
