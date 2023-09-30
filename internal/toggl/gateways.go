package toggl

import (
	"time"
)

type ProjectGateway interface {
	Get(filter *GetProjectsOpts) (TogglProjects, error)
	Create(projectName string, clientID int64) error
}

type TimeEntryGateway interface {
	Get(start, end time.Time) (TogglTimeEntries, error)
	Add(timeEntry *NewTogglTimeEntry) error
}

type ClientGateway interface {
	Get(filter *GetClientOpts) (TogglClients, error)
}
