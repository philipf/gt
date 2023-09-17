package togglservices

import "time"

type TogglTimeEntries []*TogglTimeEntry

type TogglTimeEntry struct {
	ID              int64       `json:"id"`
	WorkspaceID     int64       `json:"workspace_id"`
	ProjectID       int64       `json:"project_id"`
	TaskID          int64       `json:"task_id"`
	Billable        bool        `json:"billable"`
	Start           time.Time   `json:"start"`
	Stop            time.Time   `json:"stop"`
	Duration        int64       `json:"duration"`
	Description     string      `json:"description"`
	ServerDeletedAt interface{} `json:"server_deleted_at"`

	// Fields populated by IncludeMissingV9Fields
	Client   string
	ClientID int64
	Project  string
}
