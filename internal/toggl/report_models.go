package toggl

import (
	"fmt"
	"time"
)

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

func (r *Report) Duration() float64 {
	var totalDuration float64
	for _, projectSection := range r.ProjectSections {
		totalDuration += projectSection.Duration()
	}
	return totalDuration
}

func NewReport(sd, ed time.Time) (*Report, error) {
	// Validate start and end dates
	if sd.After(ed) {
		return nil, fmt.Errorf("start date must be before end date")
	}

	return &Report{
		StartDate:       sd,
		EndDate:         ed,
		ProjectSections: make([]*ProjectSection, 0),
	}, nil
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

func (p *ProjectSection) Duration() float64 {
	var totalDuration float64
	for _, daySection := range p.DaySections {
		totalDuration += daySection.Duration()
	}
	return totalDuration
}

func NewProjectSection() *ProjectSection {
	return &ProjectSection{
		DaySections: make([]*DaySection, 0),
	}
}

type DaySection struct {
	Date             time.Time
	IntervalSections []*IntervalSection
}

func (d *DaySection) Duration() float64 {
	var totalDuration float64
	for _, interval := range d.IntervalSections {
		totalDuration += interval.Duration
	}
	return totalDuration
}

func NewDaySection() *DaySection {
	return &DaySection{
		IntervalSections: make([]*IntervalSection, 0),
	}
}

type IntervalSection struct {
	StartDateTime time.Time
	EndDateTime   time.Time
	Description   string
	Duration      float64
}

func (i *IntervalSection) DurationS() string {
	t := time.Duration(i.Duration * float64(time.Hour))
	hours := int(t.Hours())
	minutes := int(t.Minutes()) % 60 // use modulo operation to get remaining minutes
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
