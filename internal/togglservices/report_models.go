package togglservices

import "time"

type ReportRange int

const (
	Today ReportRange = iota
	Yesterday
	LastWeek
	ThisWeek
	Custom
	EYesterday
)

type Report struct {
	StartDate       time.Time
	EndDate         time.Time
	ProjectSections []*ProjectSection
}

type ProjectKey struct {
	Client    string
	Project   string
	ProjectID int64
}

type ProjectSection struct {
	TaskID      int
	IsTask      bool
	ClientID    string
	ProjectID   string
	Client      string
	Project     string
	DaySections []*DaySection
}

type DaySection struct {
	Date             time.Time
	IntervalSections []*IntervalSection
}

type IntervalSection struct {
	StartDateTime time.Time
	EndDateTime   time.Time
	Description   string
	Duration      float64
}
