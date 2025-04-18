package gtd

import (
	"os"
	"time"

	"github.com/philipf/gt/internal/settings"
)

func AddDescriptionNote(action *Action) error {
	inTemplate := settings.GetInTemplatePath()
	if _, err := os.Stat(inTemplate); os.IsNotExist(err) {
		return err
	}

	kanbanPath := settings.GetKanbanBasePath()

	err := ActionToMd(action, inTemplate, kanbanPath)
	if err != nil {
		return err
	}
	return nil
}

func AddToKanban(todo string, withLink bool, due *time.Time) error {
	path := settings.GetKanbanBoardPath()
	err := InsertTodo(path, settings.GetKanbanInColumn(), todo, withLink, due)
	return err
}
