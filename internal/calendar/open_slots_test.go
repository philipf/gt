package calendar

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func buildDummyDay() Day {
	d := Day{
		ID:    uuid.New(),
		Date:  time.Date(2023, 8, 20, 0, 0, 0, 0, time.UTC),
		Start: time.Date(2023, 8, 20, 9, 0, 0, 0, time.UTC),
		End:   time.Date(2023, 8, 20, 17, 0, 0, 0, time.UTC),
	}

	offset := 0 // in minutes
	for i := 0; i < 5; i++ {
		startHour := 9 + offset/60
		startMin := offset % 60
		endHour := startHour + 1

		d.AddSegment(Segment{
			ID:          uuid.New(),
			Description: fmt.Sprintf("S%d", i),
			Start:       time.Date(2023, 8, 20, startHour, startMin, 0, 0, time.UTC),
			End:         time.Date(2023, 8, 20, endHour, startMin, 0, 0, time.UTC),
		})

		offset += 90 // 1 hour for the segment + 30 minutes gap
	}

	// print the segments to the console showing the start and end times
	// for _, s := range d.Segments {
	// 	fmt.Printf("Segment %s: %s - %s\n", s.Description, s.Start, s.End)
	// }

	return d
}

func TestSegments(t *testing.T) {
	day := buildDummyDay()
	expectedDescriptions := []string{"S0", "S1", "S2", "S3", "S4"}

	for i, s := range day.Segments {
		assert.Equal(t, expectedDescriptions[i], s.Description)
		// You can add more assertions for start and end times if needed
	}
}

func TestFindOpenSlots(t *testing.T) {
	day := buildDummyDay()
	openSlots := FindOpenSlots(&day)

	// Test that 5 open slots are returned
	assert.Equal(t, 5, len(openSlots))

	expectedOpenSlots := []OpenSlot{
		{
			Start: time.Date(2023, 8, 20, 10, 0, 0, 0, time.UTC),
			End:   time.Date(2023, 8, 20, 10, 30, 0, 0, time.UTC),
		},
		{
			Start: time.Date(2023, 8, 20, 11, 30, 0, 0, time.UTC),
			End:   time.Date(2023, 8, 20, 12, 0, 0, 0, time.UTC),
		},
		{
			Start: time.Date(2023, 8, 20, 13, 0, 0, 0, time.UTC),
			End:   time.Date(2023, 8, 20, 13, 30, 0, 0, time.UTC),
		},
		{
			Start: time.Date(2023, 8, 20, 14, 30, 0, 0, time.UTC),
			End:   time.Date(2023, 8, 20, 15, 0, 0, 0, time.UTC),
		},
		{
			Start: time.Date(2023, 8, 20, 16, 0, 0, 0, time.UTC),
			End:   time.Date(2023, 8, 20, 17, 0, 0, 0, time.UTC),
		},
	}

	assert.Equal(t, expectedOpenSlots, openSlots)
}
