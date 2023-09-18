package togglservices

import (
	"fmt"
	"sort"
	"testing"
	"time"
)

func TestFetchTimeEntries(t *testing.T) {
	t.Skip("skipping integration test")

	sd := time.Date(2023, 9, 14, 0, 0, 0, 0, time.Local)
	ed := sd.AddDate(0, 0, 1)

	entries, err := FetchTimeEntries(sd, ed)

	if err != nil {
		t.Errorf("expected  no error but got %v", err)
	}

	if entries == nil {
		t.Errorf("expected entries but got nil")
	}

	// print count of entries
	fmt.Printf("Found %d entries\n", len(entries))

	// sort entries by start time
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Start.Before(entries[j].Start)
	})

	for _, entry := range entries {
		// print start and end time in local time as well as the duration and description
		fmt.Printf("%s - %s (%d) - %s\n", entry.Start.Local().Format("2006-01-02 15:04"), entry.Stop.Local().Format("2006-01-02 15:04"), entry.Duration, entry.Description)
	}
}
