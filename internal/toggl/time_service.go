package toggl

import (
	"fmt"
	"log"
	"sort"
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

func (t *TimeService) Add(entry *NewTogglTimeEntry) error {
	return t.timeEntryGateway.Add(entry)
}

func (t *TimeService) ResumeLast() error {
	n := time.Now()
	sd := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.Local)
	ed := sd.AddDate(0, 0, 1).Add(-time.Second)
	entries, err := t.timeEntryGateway.Get(sd, ed)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return fmt.Errorf("nothing to resume, no time entries for today")
	}

	log.Println("Found", len(entries), "entries for today")

	// sort entries by start time
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Start.After(entries[j].Start)
	})

	lastEntry := entries[0]

	newEntry := NewTogglTimeEntry{
		Description: lastEntry.Description,
		ProjectID:   lastEntry.ProjectID,
		WorkspaceID: lastEntry.WorkspaceID,
		CreatedWith: "gt",
		Start:       time.Now().Format(time.RFC3339),
		Stop:        "", // empty stop time means the entry is still running
		Duration:    -1,
	}

	err = t.timeEntryGateway.Add(&newEntry)
	if err != nil {
		return err
	}

	return nil
}

func (t *TimeService) Stop() error {
	current, err := t.timeEntryGateway.GetCurrent()
	if err != nil {
		return err
	}

	if current == nil {
		//nothing to stop
		return fmt.Errorf("no running time entry to stop")
	}

	return t.timeEntryGateway.Stop(current.ID)
}

func (t *TimeService) EditDesc(desc string) error {
	// Get the current time timeEntryGateway

	current, err := t.timeEntryGateway.GetCurrent()
	if err != nil {
		return err
	}

	if current.ID == 0 {
		//nothing to edit
		return fmt.Errorf("no time entry to edit")
	}

	return t.timeEntryGateway.EditDesc(current.ID, desc)
}
