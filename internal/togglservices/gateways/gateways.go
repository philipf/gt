package gateways

import (
	"time"

	"github.com/philipf/gt/internal/togglservices"
)

type ClientGateway interface {
	GetClients(filter string) (togglservices.TogglClients, error)
}

type ProjectGateway interface {
	GetProjects() (togglservices.TogglProjects, error)
}

type TimeEntryGateway interface {
	GetTimeEntries(start, end time.Time) (togglservices.TogglTimeEntries, error)
}
