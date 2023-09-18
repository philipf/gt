package toggl

import (
	"time"
)

type TimeService interface {
	GetTimeEntries(start, end time.Time) (TogglTimeEntries, error)
}

type TimeServiceImplementation struct {
	TimeEntryGateway TimeEntryGateway
	ClientService    ClientService
	ProjectService   ProjectService
}

func (t *TimeServiceImplementation) GetTimeEntries(start, end time.Time) (TogglTimeEntries, error) {
	entries, err := t.TimeEntryGateway.GetTimeEntries(start, end)
	if err != nil {
		return nil, err
	}

	err = t.includeProjectAndClient(entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (t *TimeServiceImplementation) includeProjectAndClient(timeEntries TogglTimeEntries) error {
	projects, err := t.ProjectService.GetProjects(nil)
	if err != nil {
		return err
	}

	clients, err := t.ClientService.GetClients("")
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
