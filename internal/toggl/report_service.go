package toggl

import (
	"fmt"
	"sort"
	"time"
)

type ReportService struct {
	timeService    TimeService
	projectService ProjectService
}

func NewReportService(timeService TimeService, projectService ProjectService) ReportService {
	return ReportService{
		timeService:    timeService,
		projectService: projectService,
	}
}

func (r *ReportService) GetReport(sd, ed time.Time) (*Report, error) {

	timeEntries, err := r.timeService.GetTimeEntries(sd, ed)
	if err != nil {
		return nil, err
	}

	// Create report
	rpt, err := NewReport(sd, ed)
	if err != nil {
		return nil, err
	}

	// Group by ProjectKey
	projectGroups := make(map[ProjectKey]TogglTimeEntries)
	for _, entry := range timeEntries {
		key := ProjectKey{
			Client:    entry.Client,
			Project:   entry.Project,
			ProjectID: entry.ProjectID,
		}
		projectGroups[key] = append(projectGroups[key], entry)
	}

	// Extract the keys from the projectGroups map into a slice
	var keys []ProjectKey
	for k := range projectGroups {
		keys = append(keys, k)
	}

	// Sort the keys slice based on your criteria
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Client != keys[j].Client {
			return keys[i].Client < keys[j].Client
		}
		return keys[i].Project < keys[j].Project
	})

	// Sort and process each group
	for _, key := range keys {
		entries := projectGroups[key]
		titleParseResult, err := r.projectService.ParseProjectTitle(key.Project)
		if err != nil {
			fmt.Println("Error parsing title:", err)
			continue

		}

		ps := &ProjectSection{
			Client:      key.Client,
			Project:     key.Project,
			ProjectID:   fmt.Sprintf("%d", key.ProjectID),
			IsTask:      titleParseResult.IsTask,
			TaskID:      titleParseResult.TaskID,
			DaySections: make([]*DaySection, 0),
		}

		rpt.ProjectSections = append(rpt.ProjectSections, ps)

		// Group by day
		dayGroups := make(map[time.Time]TogglTimeEntries)
		for _, entry := range entries {
			day := time.Date(entry.Start.Year(), entry.Start.Month(), entry.Start.Day(), 0, 0, 0, 0, time.Local)
			dayGroups[day] = append(dayGroups[day], entry)
		}

		// Sort and process each day group
		for day, dayEntries := range dayGroups {
			ds := &DaySection{
				Date:             day,
				IntervalSections: make([]*IntervalSection, 0),
			}
			ps.DaySections = append(ps.DaySections, ds)

			// Sort entries by start time
			sort.Slice(dayEntries, func(i, j int) bool {
				return dayEntries[i].Start.Before(dayEntries[j].Start)
			})

			for _, entry := range dayEntries {
				durationHours := float64(entry.Duration) / 3600.0 // Convert seconds to hours

				// if the entry.Stop is a zero time, then use the current time. This is because the entry is still running
				if entry.Stop.IsZero() {
					entry.Stop = time.Now()
				}

				is := &IntervalSection{
					StartDateTime: entry.Start,
					EndDateTime:   entry.Stop,
					Description:   entry.Description,
					Duration:      durationHours,
				}
				ds.IntervalSections = append(ds.IntervalSections, is)
			}
		}
	}

	return rpt, nil

}
