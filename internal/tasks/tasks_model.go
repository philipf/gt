package tasks

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID           uuid.UUID
	ExternalID   string
	Title        string
	Description  string
	ExternalLink string
	CreatedAt    time.Time
	ModifiedAt   time.Time
}

func (t *Task) Create(
	externalId, title, description, externalLink string,
	createdAt, modifiedAt time.Time,
) (*Task, error) {

	return &Task{
		ID:           uuid.New(),
		ExternalID:   externalId,
		Title:        title,
		Description:  description,
		ExternalLink: externalLink,
	}, nil
}
