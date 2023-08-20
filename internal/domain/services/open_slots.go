// This (domain service) finds the open slots for a given day.
// It's a simple algorithm that iterates over the segments and finds the open slots. It's a good example of a domain service because it's not a behavior of the Day itself,
// but it's a behavior that is needed by the application. It's also a good example of a domain service because it's a behavior that is not related to a single aggregate.

package services

import (
	"fmt"
	"time"

	domain "github.com/philipf/gt/internal/domain/models"
)

type OpenSlot struct {
	Start time.Time
	End   time.Time
}

type OpenSlotsService struct {
}

func FindOpenSlots(day domain.Day) []OpenSlot {
	openSlots := []OpenSlot{}
	current := day.Start
	for _, segment := range day.Segments {
		if current.Before(segment.Start) {
			openSlots = append(openSlots, OpenSlot{
				Start: current,
				End:   segment.Start,
			})
		}
		current = segment.End
	}
	if current.Before(day.End) {
		openSlots = append(openSlots, OpenSlot{
			Start: current,
			End:   day.End,
		})
	}
	return openSlots
}

// PrintOpenSlots prints the open slots for a given day.
func PrintOpenSlots(slots []OpenSlot) {
	for _, slot := range slots {
		fmt.Printf("Slot: %s - %s\n", slot.Start.Format(time.Kitchen), slot.End.Format(time.Kitchen))
	}
}
