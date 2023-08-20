// The main entry point for the gt cli command
package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/philipf/gt/internal/domain/models"
	"github.com/philipf/gt/internal/domain/services"
)

func main() {
	fmt.Println("GT Version 0.0.1")

	// Build some dummy data
	// First, create a day with a start and end time
	d := models.Day{
		Id:    uuid.New(),
		Date:  time.Date(2023, 8, 20, 0, 0, 0, 0, time.UTC),
		Start: time.Date(2023, 8, 20, 9, 0, 0, 0, time.UTC),
		End:   time.Date(2023, 8, 20, 17, 0, 0, 0, time.UTC),
	}

	// Add 5 segments, each should last for duration of 1 hour each, with 30 minutes gaps between the start and end time of the previous segment
	offset := 0 // in minutes
	for i := 0; i < 5; i++ {
		startHour := 9 + offset/60
		startMin := offset % 60
		endHour := startHour + 1

		d.AddSegment(models.Segment{
			Id:          uuid.New(),
			Description: fmt.Sprintf("S%d", i),
			Start:       time.Date(2023, 8, 20, startHour, startMin, 0, 0, time.UTC),
			End:         time.Date(2023, 8, 20, endHour, startMin, 0, 0, time.UTC),
		})

		offset += 90 // 1 hour for the segment + 30 minutes gap
	}

	// loop through the segments and print them, showing the start and end times
	for _, s := range d.Segments {
		fmt.Printf("Segment %s: %s - %s\n", s.Description, s.Start, s.End)
	}

	// get the open slots for the day
	openSlots := services.FindOpenSlots(d)

	// loop through the open slots and print them, showing the start and end times
	for _, s := range openSlots {
		fmt.Printf("Open Slot: %s - %s\n", s.Start, s.End)
	}
}
