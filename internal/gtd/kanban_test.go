package gtd

import (
	"testing"
)

func TestFindLastToDo(t *testing.T) {
	lines := []string{
		"## In",
		"- [ ] Test 1",
		"- [ ] Test 2",
		"## Archive",
		"- [ ] Test 3",
		"- [ ] Test 4",
	}

	lastTodo := findLastToDo(lines, "Archive")

	if lastTodo != 5 {
		t.Errorf("Expected 5, got %d", lastTodo)
	}
}

func TestFindLastToDo_NoToDosForIn(t *testing.T) {
	lines := []string{
		"## In",

		"## Archive",
		"- [ ] Test 3",
		"- [ ] Test 4",
	}

	lastTodo := findLastToDo(lines, "In")

	if lastTodo != 0 {
		t.Errorf("Expected 0, got %d", lastTodo)
	}
}

func TestFindLastToDo_None(t *testing.T) {
	lines := []string{
		"",
	}

	lastTodo := findLastToDo(lines, "In")

	if lastTodo != -1 {
		t.Errorf("Expected -1, got %d", lastTodo)
	}
}

func TestFindLastToDo_OnlyOneHeadingWithTodos(t *testing.T) {
	lines := []string{
		"## In",
		"- [ ] Test 1",
		"- [ ] Test 2",
	}

	lastTodo := findLastToDo(lines, "In")

	if lastTodo != 2 {
		t.Errorf("Expected 2, got %d", lastTodo)
	}
}

func TestFindLastToDo_DoneAndExtraSpace(t *testing.T) {
	lines := []string{
		"## In",
		" - [x] Done",

		"## Archive",
		"- [ ] Test 3",
		"- [ ] Test 4",
	}

	lastTodo := findLastToDo(lines, "In")

	if lastTodo != 1 {
		t.Errorf("Expected 1, got %d", lastTodo)
	}
}
