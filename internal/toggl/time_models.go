package toggl

import "time"

type GetTimeOpts struct {
	Start time.Time
	End   time.Time
}

type TogglTimeEntries []*TogglTimeEntry

// Online documentation: https://developers.track.toggl.com/docs/api/time_entries/index.html#response-1
type TogglTimeEntry struct {
	ID              int64       `json:"id"`
	WorkspaceID     int64       `json:"workspace_id"`
	ProjectID       int64       `json:"project_id"`
	TaskID          int64       `json:"task_id"`
	Billable        bool        `json:"billable"`
	Start           time.Time   `json:"start"`
	Stop            time.Time   `json:"stop"`
	Duration        int64       `json:"duration"` // in seconds, running entries contain a negative duration
	Description     string      `json:"description"`
	ServerDeletedAt interface{} `json:"server_deleted_at"`

	// derived fields
	Client   string
	ClientID int64
	Project  string
}

// https://developers.track.toggl.com/docs/api/time_entries#post-timeentries
type NewTogglTimeEntry struct {
	Billable    bool     `json:"billable,omitempty"`    // Whether the time entry is marked as billable. Optional, default false.
	CreatedWith string   `json:"created_with"`          // Must be provided when creating a time entry and should identify the service/application used to create it.
	Description string   `json:"description,omitempty"` // Time entry description. Optional.
	Duration    int64    `json:"duration"`              // Time entry duration. For running entries should be negative, preferable -1.
	ProjectID   int64    `json:"project_id,omitempty"`  // Project ID. Optional.
	Start       string   `json:"start"`                 // Start time in UTC, required for creation. Format: 2006-01-02T15:04:05Z.
	Stop        string   `json:"stop,omitempty"`        // Stop time in UTC, can be omitted if it's still running or created with "duration". If "stop" and "duration" are provided, values must be consistent (start + duration == stop).
	TagAction   string   `json:"tag_action,omitempty"`  // Can be "add" or "delete". Used when updating an existing time entry.
	TagIDs      []int    `json:"tag_ids,omitempty"`     // IDs of tags to add/remove.
	Tags        []string `json:"tags,omitempty"`        // Names of tags to add/remove. If name does not exist as tag, one will be created automatically.
	WorkspaceID int64    `json:"workspace_id"`          // Workspace ID. Required.
}

type UpdateTogglTimeEntryDesc struct {
	Description string `json:"description"`
}
