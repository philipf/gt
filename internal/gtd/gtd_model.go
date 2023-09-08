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
	UpdatedAt    time.Time
	Due          time.Time
	Status       string
	Channel      string
	Priority     int
	// Next action
}

func CreateBasicAction(title, description, channel string) (*Action, error) {
	return CreateAction("", title, description, "", channel, time.Now(), time.Now(), In)
}

func CreateAction(
	externalID, title, description, externalLink, channel string,
	createdAt, updatedAt time.Time,
	status string,
) (*Action, error) {

	return &Action{
		ID:           uuid.New(),
		ExternalID:   externalID,
		Title:        title,
		Description:  description,
		ExternalLink: externalLink,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Status:       status,
		Channel:      channel,
	}, nil
}

// According to the GTD methodology, a project is a task that requires more than one action to complete.
// A project has specific end date and time.
type Project struct{}
