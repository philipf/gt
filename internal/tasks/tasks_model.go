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
	UpdatedAt    time.Time
	DueAt        time.Time
}

func CreateTask(
	externalID, title, description, externalLink string,
	createdAt, modifiedAt time.Time,
) (*Task, error) {

	return &Task{
		ID:           uuid.New(),
		ExternalID:   externalID,
		Title:        title,
		Description:  description,
		ExternalLink: externalLink,
		CreatedAt:    createdAt,
		UpdatedAt:    modifiedAt,
	}, nil
}
