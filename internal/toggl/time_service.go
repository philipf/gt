package toggl

import (
	"time"
)

type TimeService struct {
	timeEntryGateway TimeEntryGateway
	clientService    ClientService
	projectService   ProjectService
}

func NewTimeService(timeEntryGateway TimeEntryGateway, clientService ClientService, projectService ProjectService) TimeService {
	return TimeService{
		timeEntryGateway: timeEntryGateway,
		clientService:    clientService,
		projectService:   projectService,
	}
}

func (t *TimeService) GetTimeEntries(start, end time.Time) (TogglTimeEntries, error) {
	entries, err := t.timeEntryGateway.GetTimeEntries(start, end)
	if err != nil {
		return nil, err
	}

	err = t.includeProjectAndClient(entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (t *TimeService) includeProjectAndClient(timeEntries TogglTimeEntries) error {
	projects, err := t.projectService.GetProjects(nil)
	if err != nil {
		return err
	}

	clients, err := t.clientService.GetClients("")
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
