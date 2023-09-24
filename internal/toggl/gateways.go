package toggl

import (
	"time"
)

type ProjectGateway interface {
	GetProjects(filter *GetProjectsOpts) (TogglProjects, error)
	CreateProject(projectName string, clientID int64) error
}

type TimeEntryGateway interface {
	GetTimeEntries(start, end time.Time) (TogglTimeEntries, error)
}

type ClientGateway interface {
	GetClients(filter string) (TogglClients, error)
}

// func NewTogglClientGateway() ClientGateway {
// 	return &gateways.TogglClientGateway{}
// }

// func NewFileClientGateway() ClientGateway {
// 	return &gateways.FileClientGateway{}
// }
