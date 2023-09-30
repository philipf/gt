package toggl

type GetProjectsOpts struct {
	Name            string
	IncludeArchived bool
	ClientID        int64
}

type ProjectTitle struct {
	Client   string
	Project  string
	IsTask   bool
	TaskID   int
	TicketID string
}

type TogglProjects []*TogglProjectElement

type TogglProjectElement struct {
	ID                  int64       `json:"id"`
	WorkspaceID         int64       `json:"workspace_id"`
	ClientID            int64       `json:"client_id"`
	Name                string      `json:"name"`
	IsPrivate           bool        `json:"is_private"`
	Active              bool        `json:"active"`
	At                  string      `json:"at"`
	CreatedAt           string      `json:"created_at"`
	ServerDeletedAt     interface{} `json:"server_deleted_at"`
	Color               string      `json:"color"`
	Billable            interface{} `json:"billable"`
	Template            interface{} `json:"template"`
	AutoEstimates       interface{} `json:"auto_estimates"`
	EstimatedHours      interface{} `json:"estimated_hours"`
	Rate                interface{} `json:"rate"`
	RateLastUpdated     interface{} `json:"rate_last_updated"`
	Currency            interface{} `json:"currency"`
	Recurring           bool        `json:"recurring"`
	RecurringParameters interface{} `json:"recurring_parameters"`
	CurrentPeriod       interface{} `json:"current_period"`
	FixedFee            interface{} `json:"fixed_fee"`
	ActualHours         int64       `json:"actual_hours"`

	// Derived
	Client string
}

// ProjectsByName implements sort.Interface based on the Name field of TogglClientElement.
type ProjectsByName TogglProjects

func (a ProjectsByName) Len() int           { return len(a) }
func (a ProjectsByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ProjectsByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
