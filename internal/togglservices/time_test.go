package togglservices

import (
	"fmt"
	"testing"
	"time"
)

func TestFetchTimeEntries(t *testing.T) {

	sd := time.Date(2023, 9, 14, 0, 0, 0, 0, time.Local)
	ed := sd.AddDate(0, 0, 1)

	entries, err := FetchTimeEntries(sd, ed)

	if err != nil {
		t.Errorf("expected  no error but got %v", err)
	}

	if entries == nil {
		t.Errorf("expected entries but got nil")
	}

	//dereference the pointer and loop through the entries
	for _, entry := range entries {
		fmt.Println(*entry)
	}
}
