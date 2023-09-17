package togglservices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

// Placeholder for FetchTimeEntries function
func FetchTimeEntries(startDateTime, endDateTime time.Time) (TogglTimeEntries, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v9/me/time_entries", BASE_URI), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(viper.GetString("toggl.ApiKey"), "api_token")

	// Add parameters
	q := req.URL.Query()
	q.Add("start_date", startDateTime.Format(time.RFC3339))
	q.Add("end_date", endDateTime.Format(time.RFC3339))
	req.URL.RawQuery = q.Encode()

	//fmt.Println(req.URL.RawQuery)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status: %s", resp.Status)
	}

	var timeEntries TogglTimeEntries
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&timeEntries)
	if err != nil {
		return nil, err
	}

	// Filter out entries with ServerDeletedAt set
	var filteredEntries TogglTimeEntries
	for _, entry := range timeEntries {
		if entry.ServerDeletedAt == nil {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	return timeEntries, nil
}

func includeProjectAndClient(timeEntries TogglTimeEntries) error {
	projects, err := GetProjects(nil)
	if err != nil {
		return err
	}

	clients, err := GetClients("")
	if err != nil {
		return err
	}

	// Convert projects and clients to maps for easy lookup
	projectMap := make(map[int64]*TogglProjectElement)
	clientMap := make(map[int64]*TogglClientElement)

	for _, client := range clients {
		// make a map of clients using the client id as the key
		clientMap[client.ID] = client
	}

	for _, project := range projects {
		// make a map of projects using the project id as the key
		projectMap[project.ID] = project
		projectMap[project.ID].Client = clientMap[project.ClientID].Name
	}

	// Populate missing fields in time entries
	for i := range timeEntries {
		project, ok := projectMap[timeEntries[i].ProjectID]
		if ok {
			timeEntries[i].Client = project.Client
			timeEntries[i].ClientID = project.ClientID
			timeEntries[i].Project = project.Name
		}
	}

	return nil
}
