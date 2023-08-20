package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewDayShouldAllowNewItem(t *testing.T) {
	p := Person{
		Id:    uuid.New(),
		Name:  "Test Person",
		Email: "test@test.com",
	}

	t1 := time.Date(2023, 8, 20, 0, 0, 0, 0, time.UTC)
	d, error := p.AddDay(t1)

	if error != nil {
		t.Fatalf("Error adding day: %s", error)
	}

	if d == nil {
		t.Fatalf("Day should not be nil")
	}

	if d.Id == uuid.Nil {
		t.Fatalf("Day should have an ID")
	}

	if d.Date != t1 {
		t.Fatalf("Day should have the correct date")
	}

	if d.Start != time.Date(2023, 8, 20, 9, 0, 0, 0, time.UTC) {
		t.Fatalf("Day should have the correct start time")
	}

	if d.End != time.Date(2023, 8, 20, 17, 0, 0, 0, time.UTC) {
		t.Fatalf("Day should have the correct end time")
	}

	if len(d.Segments) != 0 {
		t.Fatalf("Day should have no segments")
	}
}

// Test to ensure that a day cannot be added if it already exists
func TestNewDayShouldNotAllowDuplicateItem(t *testing.T) {
	p := Person{
		Id:    uuid.New(),
		Name:  "Test Person",
		Email: "test@test.com",
	}

	t1 := time.Date(2023, 8, 20, 0, 0, 0, 0, time.UTC)
	_, error := p.AddDay(t1)

	if error != nil {
		t.Fatalf("Error adding day: %s", error)
	}

	t2 := time.Date(2023, 8, 20, 0, 0, 0, 0, time.UTC)
	_, error = p.AddDay(t2)

	if error == nil {
		t.Fatalf("Error should have been raised")
	}

	if error.Error() != "day already exists for the provided date" {
		t.Fatalf("Error should have been raised: day already exists for the provided date")
	}

	if len(p.Days) != 1 {
		t.Fatalf("Day should not have been added")
	}
}

// Test that when multiple days are added, they are sorted by date
func TestNewDayShouldSortItems(t *testing.T) {
	p := Person{
		Id:    uuid.New(),
		Name:  "Test Person",
		Email: "test@test.com",
	}

	t1 := time.Date(2023, 8, 20, 0, 0, 0, 0, time.UTC)
	_, error := p.AddDay(t1)

	if error != nil {
		t.Fatalf("Error adding day: %s", error)
	}

	t2 := time.Date(2023, 8, 19, 0, 0, 0, 0, time.UTC)
	_, error = p.AddDay(t2)

	if error != nil {
		t.Fatalf("Error adding day: %s", error)
	}

	if len(p.Days) != 2 {
		t.Fatalf("Day should have been added")
	}

	if p.Days[0].Date != t2 {
		t.Fatalf("Day should have been sorted")
	}

	// print days to console for debugging
	// for _, d := range p.Days {
	// 	fmt.Println(d.Date)
	// }
}
