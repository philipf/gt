// Contains the domain model for the GTD (getting things done) by David Allen
package gtd

import (
	"time"

	"github.com/google/uuid"
)

// Status of an action can be
const (
	In = "In"
	Eliminated
	Incubate
	Referenced
	InProgress
	//Project
	WaitingFor
	Delegated
	Deferred
	Scheduled
)

// An action is the smallest unit of work that can be done in a single step.
type Action struct {
	ID           uuid.UUID
	ExternalID   string
	Title        string
	Description  string
	ExternalLink string
	CreatedAt    time.Time
	ModifiedAt   time.Time
	Due          time.Time
	Status       string
	Channel      string
	Priority     int
	// Next action
}

func CreateAction(
	externalID, title, description, externalLink string,
	createdAt, modifiedAt time.Time,
	status string,
) (*Action, error) {

	return &Action{
		ID:           uuid.New(),
		ExternalID:   externalID,
		Title:        title,
		Description:  description,
		ExternalLink: externalLink,
		CreatedAt:    createdAt,
		ModifiedAt:   modifiedAt,
		Status:       status,
	}, nil
}

// According to the GTD methodology, a project is a task that requires more than one action to complete.
// A project has specific end date and time.
type Project struct{}
