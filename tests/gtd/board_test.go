package gtd

import (
	"testing"

	"github.com/philipf/gt/internal/gtd"
)

func TestInsertItem(t *testing.T) {

	err := gtd.InsertTodo("_Board.md", "In", "Test 1")
	if err != nil {
		t.Error(err)
	}

}
