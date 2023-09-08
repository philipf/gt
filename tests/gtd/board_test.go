package gtd

import (
	"os"
	"testing"

	"github.com/philipf/gt/internal/files"
	"github.com/philipf/gt/internal/gtd"
)

func TestInsertItemWithLink(t *testing.T) {

	const (
		testFile = "TestInsertItemWithLink.md"
		expected = "TestInsertItemWithLink_expected.md"
	)

	files.CopyFile("_Board.md", testFile)

	err := gtd.InsertTodo(testFile, "In", "Test 1", true)
	if err != nil {
		t.Error(err)
	}

	equal, err := files.AreFilesEqual(testFile, expected)
	if err != nil {
		t.Error(err)
	}

	if !equal {
		t.Error("Files are not equal")
	}

	os.Remove(testFile)
}

func TestInsertItem(t *testing.T) {

	const (
		testFile = "TestInsertItem.md"
		expected = "TestInsertItem_expected.md"
	)

	files.CopyFile("_Board.md", testFile)

	err := gtd.InsertTodo(testFile, "In", "Test 1", false)
	if err != nil {
		t.Error(err)
	}

	equal, err := files.AreFilesEqual(testFile, expected)
	if err != nil {
		t.Error(err)
	}

	if !equal {
		t.Error("Files are not equal")
	}

	os.Remove(testFile)
}
