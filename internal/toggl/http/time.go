package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/philipf/gt/internal/cache"
	"github.com/philipf/gt/internal/settings"
	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/viper"
)

type TogglTimeEntriesGateway struct {
}

func (t *TogglTimeEntriesGateway) Get(start, end time.Time) (toggl.TogglTimeEntries, error) {
	log.Printf("Getting time entries from %s to %s\n", start, end)
	cache := cache.JsonFileCache[toggl.TogglTimeEntries, toggl.GetTimeOpts]{}

	cacheDir, err := settings.GetGtConfigPath()
	if err != nil {
		return nil, err
	}

	cacheFilePath := path.Join(cacheDir, "toggl-time-entries.json")
	cacheMaxAge := time.Duration(time.Minute * 2)

	filter := &toggl.GetTimeOpts{
		Start: start,
		End:   end,
	}

	cr, err := cache.Get(*filter, cacheFilePath, cacheMaxAge)

	if err != nil {
		log.Println("Cache miss")
	} else {
		log.Println("Cache hit")
		return *cr, nil
	}

	r, err := t.getFromToggl(filter.Start, filter.End)
	if err != nil {
		return nil, err
	}

	err = cache.Save(*filter, cacheFilePath, &r, cacheMaxAge)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (t *TogglTimeEntriesGateway) getFromToggl(start, end time.Time) (toggl.TogglTimeEntries, error) {
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
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(viper.GetString("toggl.ApiKey"), "api_token")

	log.Println("URI:", u.String())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status: %s", resp.Status)
	}

	// print response body
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// return nil, err
	// }
	// fmt.Println(string(body))

	var timeEntries toggl.TogglTimeEntries
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&timeEntries)
	if err != nil {
		return nil, err
	}

	// Filter out entries with ServerDeletedAt set
	var filteredEntries toggl.TogglTimeEntries
	for _, entry := range timeEntries {
		entry.Start = entry.Start.Local()
		entry.Stop = entry.Stop.Local()
		// fmt.Println("Entry:", entry)
		if entry.ServerDeletedAt == nil {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	return timeEntries, nil
}

func (t *TogglTimeEntriesGateway) Add(timeEntry *toggl.NewTogglTimeEntry) error {
	uri, err := getTimeEntriesWorkspaceUri()
	if err != nil {
		return err
	}

	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	q := u.Query()
	u.RawQuery = q.Encode()

	client := &http.Client{}
	payload, err := json.Marshal(timeEntry)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(getAPIToken(), "api_token")

	log.Println("URI:", u.String())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// print response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status: %s", resp.Status)
	}
	return nil
}

func (t *TogglTimeEntriesGateway) GetCurrent() (*toggl.TogglTimeEntry, error) {
	uri, err := getTimeEntriesUri()
	if err != nil {
		return nil, err
	}

	uri += "/current"

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	u.RawQuery = q.Encode()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(viper.GetString("toggl.ApiKey"), "api_token")

	log.Println("URI:", u.String())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status: %s", resp.Status)
	}

	var timeEntry toggl.TogglTimeEntry
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&timeEntry)
	if err != nil {
		return nil, err
	}

	return &timeEntry, nil
}

func (t *TogglTimeEntriesGateway) Stop(entryID int64) error {
	uri, err := getTimeEntriesStopUri(entryID)
	if err != nil {
		return err
	}

	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	q := u.Query()
	u.RawQuery = q.Encode()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPatch, u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(viper.GetString("toggl.ApiKey"), "api_token")

	log.Println("URI:", u.String())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status: %s", resp.Status)
	}

	return nil
}

func (t *TogglTimeEntriesGateway) EditDesc(entryID int64, desc string) error {
	uri, err := getTimeEntriesIDUri(entryID)
	if err != nil {
		return err
	}

	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	q := u.Query()
	u.RawQuery = q.Encode()

	timeEntry := toggl.UpdateTogglTimeEntryDesc{
		Description: desc,
	}

	client := &http.Client{}
	payload, err := json.Marshal(timeEntry)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(viper.GetString("toggl.ApiKey"), "api_token")

	log.Println("URI:", u.String())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status: %s", resp.Status)
	}

	return nil
}
