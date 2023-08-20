// Package domain defines the main constructs for managing daily work schedules.
package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	// WorkDay indicates a day where work activities occur.
	WorkDay = iota

	// NonWorkingDay indicates a day free from work activities.
	NonWorkingDay
)

// Day represents a specific date where work activities might take place.
// It includes details like start and end times of work, and is influenced by factors
// such as a person's calendar entries, public holidays, and preferred working hours.
type Day struct {
	ID uuid.UUID

	// Date specifies the calendar date for the work day, excluding the time details.
	Date time.Time

	// Start marks the beginning of the work day.
	Start time.Time

	// End signifies the conclusion of the work day.
	End time.Time

	// Segments represents chunks of time within the day. These can be continuous or
	// can overlap (e.g., due to double bookings). A day may also have no segments.
	Segments []Segment
}

// AddSegment attaches a new time segment to the current day.
func (d *Day) AddSegment(s Segment) error {
	err := validateSegment(d, s)
	if err != nil {
		return err
	}

	d.Segments = append(d.Segments, s)

	return nil
}

// validateSegment ensures the segment adheres to constraints related to start and end times.
func validateSegment(d *Day, s Segment) error {
	// Ensure segment's start time precedes its end time.
	if s.Start.After(s.End) {
		return fmt.Errorf("start time (%s) is after the end time (%s)", s.Start, s.End)
	}

	// Ensure segment's start time isn't before the day's start time.
	if s.Start.Before(d.Start) {
		return fmt.Errorf("start time (%s) is before the start of the day (%s)", s.Start, d.Start)
	}

	// Ensure segment's end time doesn't exceed the day's end time.
	if s.End.After(d.End) {
		return fmt.Errorf("end time (%s) is after the end of the day (%s)", s.End, d.End)
	}

	return nil
}

// Segment defines a time interval within a day, representing either working or non-working hours.
type Segment struct {
	ID uuid.UUID

	// Description offers more details about the segment's nature.
	Description string

	// Start indicates when the segment begins, and it should fall within the encompassing day's duration.
	Start time.Time

	// End indicates when the segment concludes, and it should also stay within the day's limits.
	End time.Time

	// IsWorkingTime denotes whether the segment represents productive work time or a break.
	IsWorkingTime bool
}

// factory method for creating a new segment
func NewSegment(description string, start time.Time, end time.Time, isWorkingTime bool) Segment {
	return Segment{
		ID:            uuid.New(),
		Description:   description,
		Start:         start,
		End:           end,
		IsWorkingTime: isWorkingTime,
	}
}
