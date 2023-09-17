package calendar

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewDayShouldAllowNewItem(t *testing.T) {
	p := Person{
		ID:    uuid.New(),
		Name:  "Test Person",
		Email: "test@test.com",
	}

	t1 := time.Date(2023, 8, 20, 0, 0, 0, 0, time.UTC)
	d, error := p.AddDay(t1)

	assert.NoError(t, error, "Error adding day")
	assert.NotNil(t, d, "Day should not be nil")
	assert.NotEqual(t, uuid.Nil, d.ID, "Day should have an ID")
	assert.Equal(t, t1, d.Date, "Day should have the correct date")
	assert.Equal(t, time.Date(2023, 8, 20, 9, 0, 0, 0, time.UTC), d.Start, "Day should have the correct start time")
	assert.Equal(t, time.Date(2023, 8, 20, 17, 0, 0, 0, time.UTC), d.End, "Day should have the correct end time")
	assert.Empty(t, d.Segments, "Day should have no segments")
}

func TestNewDayShouldNotAllowDuplicateItem(t *testing.T) {
	p := Person{
		ID:    uuid.New(),
		Name:  "Test Person",
		Email: "test@test.com",
	}

	t1 := time.Date(2023, 8, 20, 0, 0, 0, 0, time.UTC)
	_, err := p.AddDay(t1)

	assert.NoError(t, err, "Error adding day")

	t2 := time.Date(2023, 8, 20, 0, 0, 0, 0, time.UTC)
	_, err = p.AddDay(t2)

	assert.Error(t, err, "Error should have been raised")

	var errExpected ErrDayAlreadyExists
	assert.ErrorAs(t, err, &errExpected)
	assert.Equal(t, t2, errExpected.Date, t2)

	assert.Len(t, p.Days, 1, "Day should not have been added")
}

func TestNewDayShouldSortItems(t *testing.T) {
	p := Person{
		ID:    uuid.New(),
		Name:  "Test Person",
		Email: "test@test.com",
	}

	t1 := time.Date(2023, 8, 20, 0, 0, 0, 0, time.UTC)
	_, error := p.AddDay(t1)

	assert.NoError(t, error, "Error adding day")

	t2 := time.Date(2023, 8, 19, 0, 0, 0, 0, time.UTC)
	_, error = p.AddDay(t2)

	assert.NoError(t, error, "Error adding day")
	assert.Len(t, p.Days, 2, "Days should have been added")
	assert.Equal(t, t2, p.Days[0].Date, "Day should have been sorted")
}
