// Package gtd provides the domain model for the GTD (Getting Things Done) methodology by David Allen.
package gtd

import (
	"time"

	"github.com/google/uuid"
)

// Action statuses represent the various states an action can be in.
const (
	In         = "In" // Incoming task or item
	Eliminated        // Task has been removed or is not actionable
	Incubate          // Task is on hold or under consideration
	Referenced        // Task refers to external resources or references
	InProgress        // Task is currently being worked on
	WaitingFor        // Waiting for an external event or person to progress
	Delegated         // Task has been handed off to someone else
	Deferred          // Task has been postponed for a later time
	Scheduled         // Task has a specific planned start time
)

// Action represents a single step or task. It's the smallest unit of work that can be done in one step.
type Action struct {
	ID           uuid.UUID // Unique ID for the action
	ExternalID   string    // External system reference ID, if applicable
	Title        string    // Short summary or title of the action
	Description  string    // Detailed description of the action
	ExternalLink string    // URL or link to external resources related to the action
	CreatedAt    time.Time // Timestamp when the action was created
	UpdatedAt    time.Time // Timestamp when the action was last updated
	Due          time.Time // Due date and time for the action
	Status       string    // Current status of the action
	Channel      string    // The medium or platform where the task was captured or will be performed
	Priority     int       // Priority of the action; higher value indicates higher priority
	// Future fields for "next action" or dependencies can be added here
}

// CreateBasicAction initializes an action with basic fields and sets default values for timestamps and status.
func CreateBasicAction(title, description, channel string) (*Action, error) {
	return CreateAction("", title, description, "", channel, time.Now(), time.Now(), In)
}

// CreateAction initializes a new Action struct with the provided parameters.
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

// Project represents a task that requires multiple actions to complete, according to the GTD methodology.
// Unlike actions, projects typically have a clear objective and end date.
type Project struct {
	// Fields related to the project can be added here
}
