// Package calendar provides functionality related to managing
// days and time segments within those days.
package calendar

import (
	"fmt"
	"time"
)

// ErrDayAlreadyExists indicates that a day already exists for a specific date.
type ErrDayAlreadyExists struct {
	Date time.Time
}

// Error returns a string representation of ErrDayAlreadyExists.
func (e ErrDayAlreadyExists) Error() string {
	return fmt.Sprintf("day already exists for the provided date: %v", e.Date)
}

// ErrDayNotFound indicates that a day was not found for a specific date.
type ErrDayNotFound struct {
	Date time.Time
}

// Error returns a string representation of ErrDayNotFound.
func (e ErrDayNotFound) Error() string {
	return fmt.Sprintf("day not found for the provided date: %v", e.Date)
}

// ErrSegmentOutsideOfDay indicates that a time segment is outside the boundaries of a day.
// This error is used when trying to associate a segment with a day and the segment's time range
// does not fit within the day's time range.
type ErrSegmentOutsideOfDay struct {
	Segment *Segment
	Day     *Day
}

// Error checks which part of the segment is outside of the day's boundaries
// and returns the appropriate error message.
func (e ErrSegmentOutsideOfDay) Error() string {
	if e.Segment.Start.Before(e.Day.Start) {
		return fmt.Sprintf("segment (%v) starts before the day (%v)", e.Segment.Start, e.Day.Start)
	}
	return fmt.Sprintf("segment (%v) ends after the day (%v)", e.Segment.End, e.Day.End)
}

// ErrInvalidSegmentRange indicates that the start and end times of a segment are not valid.
// For example, when the start time is after the end time.
type ErrInvalidSegmentRange struct {
	Start time.Time
	End   time.Time
}

// Error returns a string representation of ErrInvalidSegmentRange.
// It checks if the start time is after the end time, if so, it provides an appropriate
// error message. Otherwise, it indicates a generic invalid range error.
func (e ErrInvalidSegmentRange) Error() string {
	// Start date after end date.
	if e.Start.After(e.End) {
		return fmt.Sprintf("start time (%s) is after the end time (%v)", e.Start, e.End)
	}
	return fmt.Sprintf("invalid segment range: (%v) - (%v)", e.Start, e.End)
}
