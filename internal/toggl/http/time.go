package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/viper"
)

type TimeEntryGateway struct {
}

func (t *TimeEntryGateway) GetTimeEntries(start, end time.Time) (toggl.TogglTimeEntries, error) {
	uri, err := getTimeEntriesUri()
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("start_date", start.Format(time.RFC3339))
	q.Add("end_date", end.Format(time.RFC3339))
	u.RawQuery = q.Encode()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(viper.GetString("toggl.ApiKey"), "api_token")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status: %s", resp.Status)
	}

	var timeEntries toggl.TogglTimeEntries
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&timeEntries)
	if err != nil {
		return nil, err
	}

	// Filter out entries with ServerDeletedAt set
	var filteredEntries toggl.TogglTimeEntries
	for _, entry := range timeEntries {
		if entry.ServerDeletedAt == nil {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	return timeEntries, nil
}
