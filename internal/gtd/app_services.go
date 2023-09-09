package gtd

import (
	"os"

	"github.com/philipf/gt/internal/settings"
)

func AddDescriptionNote(action *Action) error {
	inTemplate := settings.GetInTemplatePath()
	if _, err := os.Stat(inTemplate); os.IsNotExist(err) {
		return err
	}

	inboxPath := settings.GetKanbanInboxPath()

	err := ActionToMd(action, inTemplate, inboxPath)
	if err != nil {
		return err
	}
	return nil
}

func AddToKanban(todo string, withLink bool) error {
	path := settings.GetKanbanBoardPath()
	err := InsertTodo(path, settings.GetKanbanInColumn(), todo, withLink)
	return err
}
