// Aggregate root for a day
package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	WorkDay       = iota // A workday is a day that a person works
	NonWorkingDay        // A non-working day is a day that a person does not work
)

// A Day is a day that a person works.
// It is derived from the person's calendar entries, public holidays and the persons preferred working hours
type Day struct {
	ID uuid.UUID
	// The date of the workday, without the time component
	Date time.Time

	// The time the person started work
	Start time.Time

	// The time the person finished work
	End time.Time

	// A collection of segments that make up the day
	// It is very possible that a day:
	// - Has no segments
	// - Is partially allocated with segments
	// - Is fully allocated with segments
	// - Has overlapping segments (for example double bookings)
	Segments []Segment
}

func (d *Day) AddSegment(s Segment) error {
	// Add the segment to the day
	error := validateSegment(d, s)
	if error != nil {
		return error
	}

	d.Segments = append(d.Segments, s)

	return nil
}

func validateSegment(d *Day, s Segment) error {
	// - start time is before end time
	if s.Start.After(s.End) {
		return fmt.Errorf("start time (%s) is after the end time (%s)", s.Start, s.End)
	}

	// - start time is between start and end of the day
	if s.Start.Before(d.Start) {
		// return a formatted error
		return fmt.Errorf("start time (%s) is before the start of the day (%s)", s.Start, d.Start)
	}

	// - end time is between start and end of the day
	if s.End.After(d.End) {
		return fmt.Errorf("end time (%s) is after the end of the day (%s)", s.End, d.End)
	}

	return nil
}

// A Segment is a slot in a person's day presenting working and non-working time
type Segment struct {
	ID uuid.UUID

	// A description of the segment
	Description string

	// The time the segment started.
	// The start time should:
	// - be before the end time
	// - between the start and end of the day
	Start time.Time

	// The time the segment ended
	// The end time should:
	// - be after the start time
	// - between the start and end of the day
	End time.Time

	// A flag to indicate if the segment is working or non-working
	IsWorking bool
}
