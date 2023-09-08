// Package gtd (Getting Things Done) offers utilities to manage and convert task-related actions.
package gtd

import (
	"errors"
	"os"
	"path/filepath"
	"text/template"

	"github.com/philipf/gt/internal/files"
	"github.com/philipf/gt/internal/tasks"
)

// ActionToMd converts a single Action to a markdown file using the provided template.
// The markdown file is saved in the specified output directory.
//
// Parameters:
//   - action: The action to convert.
//   - templateFile: The path to the template file.
//   - outputPath: The directory where the markdown file will be saved.
//
// Returns:
//   - error: An error if any occurred during the conversion or file operations, otherwise nil.
func ActionToMd(action *Action, templateFile, outputPath string) error {
	return ActionsToMd(&[]Action{*action}, templateFile, outputPath)
}

// ActionsToMd converts multiple Actions to markdown files using the provided template.
// Each Action results in a separate markdown file saved in the specified output directory.
//
// Parameters:
//   - actions: The list of actions to convert.
//   - templateFile: The path to the template file.
//   - outputPath: The directory where the markdown files will be saved.
//
// Returns:
//   - error: An error if any occurred during the conversion or file operations, otherwise nil.
func ActionsToMd(actions *[]Action, templateFile, outputPath string) error {
	// Parse the given template file.
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	for _, action := range *actions {
		// Ensure the output directory exists, if not, create it.
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			os.Mkdir(outputPath, 0755)
		}

		// Generate the filename for the markdown file based on the action's title.
		filename := filepath.Join(outputPath, files.ToValidFilename(action.Title)+".md")

		// Check if a file with the generated name already exists.
		if _, err := os.Stat(filename); err == nil {
			return errors.New("file already exists: " + filename)
		}

		// Open a new file for writing the action in markdown format.
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		// Ensure the file gets closed after operations are done.
		defer file.Close()

		// Apply the template to the action and write the result to the file.
		err = tmpl.Execute(file, action)
		if err != nil {
			return err
		}
	}

	return nil
}

// MapTasks converts a list of Task objects into a list of Action objects.
//
// Parameters:
//   - tasks: The list of tasks to convert.
//
// Returns:
//   - *[]Action: A pointer to a slice of converted Action objects.
//   - error: An error if any occurred during the mapping, otherwise nil.
func MapTasks(tasks []tasks.Task) (*[]Action, error) {
	actions := make([]Action, len(tasks))

	// Map each task to an action.
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
