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

// Channel defintions for goroutines in GetTimeEntries
type clientResult struct {
	clients TogglClients
	err     error
}

type projectResult struct {
	projects TogglProjects
	err      error
}

type timeEntriesResult struct {
	timeEntries TogglTimeEntries
	err         error
}

func (t *TimeService) Get(start, end time.Time) (TogglTimeEntries, error) {
	chTimeEntries := make(chan timeEntriesResult)
	chClients := make(chan clientResult)
	chProjects := make(chan projectResult)

	// Fetch time entries, clients and projects in parallel
	go func() {
		entries, err := t.timeEntryGateway.Get(start, end)
		chTimeEntries <- timeEntriesResult{entries, err}
	}()

	go func() {
		clients, err := t.clientService.Get(nil)
		chClients <- clientResult{clients, err}
	}()

	go func() {
		projects, err := t.projectService.Get(nil)
		chProjects <- projectResult{projects, err}
	}()

	timeEntriesResult := <-chTimeEntries
	if timeEntriesResult.err != nil {
		return nil, timeEntriesResult.err
	}

	clientsResult := <-chClients
	if clientsResult.err != nil {
		return nil, clientsResult.err
	}

	projectsResult := <-chProjects
	if projectsResult.err != nil {
		return nil, projectsResult.err
	}

	err := t.includeProjectAndClient(timeEntriesResult.timeEntries, clientsResult.clients, projectsResult.projects)
	if err != nil {
		return nil, err
	}

	return timeEntriesResult.timeEntries, nil
}

func (t *TimeService) includeProjectAndClient(timeEntries TogglTimeEntries, clients TogglClients, projects TogglProjects) error {
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

func (t *TimeService) Add(entry *TogglTimeEntry) error {
	//return t.timeEntryGateway.Add(entry)
	return nil
}
