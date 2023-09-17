// A person is the entity for which we associate their Calendars, Working days, Holidays and preferred working hours.
package calendar

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID    uuid.UUID
	Name  string
	Email string

	// This is the ID of the person in the external system (e.g. Google, Azure AD).
	ExternalIDs []ExternalID

	// This is a flag to indicate if this is the person who is currently logged in.
	IsSignedOnUser bool

	// Collection of days associated with this person.
	Days []Day
}

type ExternalID struct {
	ID   string
	Type string
}

// Factory method to create a new day for a person
// The following rules apply:
// - The day must not already exist for the person.
// - Only the date is used to create the day.
// - The start and end times are set to the default values.
// - The day is initialised with no segments.
func (p *Person) AddDay(date time.Time) (*Day, error) {
	// time the date to midnight.
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	// Check if the day already exists.
	for _, d := range p.Days {
		if d.Date.Equal(date) {
			return nil, ErrDayAlreadyExists{Date: date}
		}
	}

	// Create a new day.
	day := &Day{
		ID:   uuid.New(),
		Date: date,

		// TODO: remove defaults
		// Default start date for a day is 9am.
		Start: time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, date.Location()),

		// TODO: remove defaults
		// Default end date for a day is 5pm.
		End: time.Date(date.Year(), date.Month(), date.Day(), 17, 0, 0, 0, date.Location()),
	}

	// validate the day.
	if err := day.validateDay(); err != nil {
		return nil, err
	}

	// If this is the first day, then we can just append it to the slice.
	if len(p.Days) == 0 {
		p.Days = append(p.Days, *day)
		return day, nil
	} else {
		// Insert the day into the correct chronological position in the slice.
		for i, d := range p.Days {
			// The condition checks if the date of the current Day in the slice
			// (d.Date) is after the date of the new Day we're trying to insert.
			// This implies that our new Day should be inserted before this current Day
			// to keep the slice sorted by date.
			if d.Date.After(date) {
				// 1. p.Days[:i] takes a slice of all Days up to, but not including, the current Day.
				// 2. append([]Day{*day} adds our new Day to a temporary slice.
				// 3. p.Days[i:]... takes a slice of all Days from the current Day to the end.
				// 4. The outermost append joins the slices from 1 & 2 and 3 together.

				p.Days = append(p.Days[:i], append([]Day{*day}, p.Days[i:]...)...)

				// Once the new Day is inserted, we can return it and exit the function.
				return day, nil
			}
		}
	}

	return day, nil
}

// RemoveDay removes the day from the person.
// The following rules apply:
// - The day must exist for the person.
func (p *Person) RemoveDay(date time.Time) error {
	// time the date to midnight.
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	// Check if the day already exists.
	for i, d := range p.Days {
		if d.Date.Equal(date) {
			p.Days = append(p.Days[:i], p.Days[i+1:]...)
			return nil
		}
	}

	return ErrDayNotFound{Date: date}
}
