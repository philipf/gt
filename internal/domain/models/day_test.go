// Unit tests for the Day domain object.
// It is testing a number of scenarios in the internal domain package.
package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
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

	// Add a segment that covers the whole day
	s1 := Segment{
		Id:          uuid.New(),
		Description: "S1",
		Start:       time.Date(2023, 8, 20, 9, 0, 0, 0, time.UTC),
		End:         time.Date(2023, 8, 20, 17, 0, 0, 0, time.UTC),
	}

	error := d.AddSegment(s1)

	if error != nil {
		t.Fatalf("Error adding segment: %s", error)
	}
}

func TestShouldAllowMultipleSegments(t *testing.T) {
	d := initDay()

	error := d.AddSegment(Segment{
		Id:          uuid.New(),
		Description: "S1",
		Start:       time.Date(2023, 8, 20, 9, 0, 0, 0, time.UTC),
		End:         time.Date(2023, 8, 20, 10, 0, 0, 0, time.UTC),
	})

	if error != nil {
		t.Fatalf("Error adding segment: %s", error)
	}

	error = d.AddSegment(Segment{
		Id:          uuid.New(),
		Description: "S2",
		Start:       time.Date(2023, 8, 20, 11, 0, 0, 0, time.UTC),
		End:         time.Date(2023, 8, 20, 12, 0, 0, 0, time.UTC),
	})

	if error != nil {
		t.Fatalf("Error adding segment: %s", error)
	}

	// count the number of segments and raise an error if it's not 2
	if len(d.Segments) != 2 {
		t.Fatalf("Expected 2 segments, got %d", len(d.Segments))
	}
}

func TestShouldFailSegmentExceedsEndTime(t *testing.T) {
	d := initDay()

	// Add a segment that covers the whole day and then some
	s1 := Segment{
		Id:          uuid.New(),
		Description: "S1",
		Start:       time.Date(2023, 8, 20, 9, 0, 0, 0, time.UTC),
		End:         time.Date(2023, 8, 20, 17, 1, 0, 0, time.UTC),
	}

	error := d.AddSegment(s1)

	if error == nil {
		t.Fatalf("Should have resulted in an error")
	}

	if error.Error() != "end time (2023-08-20 17:01:00 +0000 UTC) is after the end of the day (2023-08-20 17:00:00 +0000 UTC)" {
		t.Fatalf("Unexpected error: %s", error)
	}
}

func TestShouldFailSegmentExceedsStartTime(t *testing.T) {
	d := initDay()

	// Add a segment that covers the whole day but starts early
	s1 := Segment{
		Id:          uuid.New(),
		Description: "S1",
		Start:       time.Date(2023, 8, 20, 8, 59, 0, 0, time.UTC),
		End:         time.Date(2023, 8, 20, 17, 0, 0, 0, time.UTC),
	}

	error := d.AddSegment(s1)

	if error == nil {
		t.Fatalf("Should have resulted in an error")
	}

	if error.Error() != "start time (2023-08-20 08:59:00 +0000 UTC) is before the start of the day (2023-08-20 09:00:00 +0000 UTC)" {
		t.Fatalf("Unexpected error: %s", error)
	}
}

func TestShouldFailSegmentEndTimeBeforeStart(t *testing.T) {
	d := initDay()

	// Add a segment that ends before it starts
	s1 := Segment{
		Id:          uuid.New(),
		Description: "S1",
		Start:       time.Date(2023, 8, 20, 11, 00, 0, 0, time.UTC), //11:00 am
		End:         time.Date(2023, 8, 20, 10, 0, 0, 0, time.UTC),  //10:00 am, this is before the start time and invalid
	}

	error := d.AddSegment(s1)

	if error == nil {
		t.Fatalf("Should have resulted in an error")
	}

	if error.Error() != "start time (2023-08-20 11:00:00 +0000 UTC) is after the end time (2023-08-20 10:00:00 +0000 UTC)" {
		t.Fatalf("Unexpected error: %s", error)
	}
}

func TestClearSegments(t *testing.T) {
	d := initDay()

	// Add a segment that covers the whole day
	s1 := Segment{
		Id:          uuid.New(),
		Description: "S1",
		Start:       time.Date(2023, 8, 20, 9, 0, 0, 0, time.UTC),
		End:         time.Date(2023, 8, 20, 17, 0, 0, 0, time.UTC),
	}

	error := d.AddSegment(s1)

	if error != nil {
		t.Fatalf("Error adding segment: %s", error)
	}

	// Clear the segments
	d.ClearSegments()

	// count the number of segments and raise an error if it's not 0
	if len(d.Segments) != 0 {
		t.Fatalf("Expected 0 segments, got %d", len(d.Segments))
	}

	// ensure add segment works again
	error = d.AddSegment(s1)

	if error != nil {
		t.Fatalf("Error adding segment: %s", error)
	}

	// count the number of segments and raise an error if it's not 1
	if len(d.Segments) != 1 {
		t.Fatalf("Expected 1 segments, got %d", len(d.Segments))
	}
}
