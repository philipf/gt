package calendar

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func initDay() Day {
	d := Day{
		Id:    uuid.New(),
		Date:  time.Date(2023, 8, 20, 0, 0, 0, 0, time.UTC),
		Start: time.Date(2023, 8, 20, 9, 0, 0, 0, time.UTC),
		End:   time.Date(2023, 8, 20, 17, 0, 0, 0, time.UTC),
	}
	return d
}

func TestShouldAllowFullSegment(t *testing.T) {
	d := initDay()
	s1 := Segment{
		Id:          uuid.New(),
		Description: "S1",
		Start:       time.Date(2023, 8, 20, 9, 0, 0, 0, time.UTC),
		End:         time.Date(2023, 8, 20, 17, 0, 0, 0, time.UTC),
	}
	assert.NoError(t, d.AddSegment(s1), "Error adding segment")
}

func TestShouldAllowMultipleSegments(t *testing.T) {
	d := initDay()

	err := d.AddSegment(Segment{
		Id:          uuid.New(),
		Description: "S1",
		Start:       time.Date(2023, 8, 20, 9, 0, 0, 0, time.UTC),
		End:         time.Date(2023, 8, 20, 10, 0, 0, 0, time.UTC),
	})
	assert.NoError(t, err, "Error adding segment S1")

	err = d.AddSegment(Segment{
		Id:          uuid.New(),
		Description: "S2",
		Start:       time.Date(2023, 8, 20, 11, 0, 0, 0, time.UTC),
		End:         time.Date(2023, 8, 20, 12, 0, 0, 0, time.UTC),
	})
	assert.NoError(t, err, "Error adding segment S2")

	assert.Equal(t, 2, len(d.Segments), "Expected 2 segments")
}

func TestShouldFailSegmentExceedsEndTime(t *testing.T) {
	d := initDay()

	s1 := Segment{
		Id:          uuid.New(),
		Description: "S1",
		Start:       time.Date(2023, 8, 20, 9, 0, 0, 0, time.UTC),
		End:         time.Date(2023, 8, 20, 17, 1, 0, 0, time.UTC),
	}

	err := d.AddSegment(s1)

	expectedErrorMsg := "segment (2023-08-20 17:01:00 +0000 UTC) ends after the day (2023-08-20 17:00:00 +0000 UTC)"
	assert.EqualError(t, err, expectedErrorMsg, "Unexpected error")
}

func TestShouldFailSegmentExceedsStartTime(t *testing.T) {
	d := initDay()

	s1 := Segment{
		Id:          uuid.New(),
		Description: "S1",
		Start:       time.Date(2023, 8, 20, 8, 59, 0, 0, time.UTC),
		End:         time.Date(2023, 8, 20, 17, 0, 0, 0, time.UTC),
	}

	err := d.AddSegment(s1)
	expectedErrorMsg := "segment (2023-08-20 08:59:00 +0000 UTC) starts before the day (2023-08-20 09:00:00 +0000 UTC)"
	assert.EqualError(t, err, expectedErrorMsg, "Unexpected error")
}

func TestShouldFailSegmentEndTimeBeforeStart(t *testing.T) {
	d := initDay()

	s1 := Segment{
		Id:          uuid.New(),
		Description: "S1",
		Start:       time.Date(2023, 8, 20, 11, 00, 0, 0, time.UTC),
		End:         time.Date(2023, 8, 20, 10, 0, 0, 0, time.UTC),
	}

	err := d.AddSegment(s1)

	var expectedError ErrInvalidSegmentRange
	assert.ErrorAs(t, err, &expectedError)
	assert.Equal(t, expectedError.Start, s1.Start)
	assert.Equal(t, expectedError.End, s1.End)

	expectedErrorMsg := "start time (2023-08-20 11:00:00 +0000 UTC) is after the end time (2023-08-20 10:00:00 +0000 UTC)"
	assert.EqualError(t, err, expectedErrorMsg, "Unexpected error")
}

func TestClearSegments(t *testing.T) {
	d := initDay()

	s1 := Segment{
		Id:          uuid.New(),
		Description: "S1",
		Start:       time.Date(2023, 8, 20, 9, 0, 0, 0, time.UTC),
		End:         time.Date(2023, 8, 20, 17, 0, 0, 0, time.UTC),
	}

	assert.NoError(t, d.AddSegment(s1), "Error adding segment")

	d.ClearSegments()

	assert.Equal(t, 0, len(d.Segments), "Expected 0 segments after clearing")

	err := d.AddSegment(s1)
	assert.NoError(t, err, "Error adding segment after clearing")

	assert.Equal(t, 1, len(d.Segments), "Expected 1 segment after adding")
}
