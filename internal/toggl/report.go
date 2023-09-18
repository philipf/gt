package toggl

// import (
// 	"fmt"
// 	"sort"
// 	"time"
// )

// func GetReportRange(args []string) ReportRange {
// 	if len(args) > 0 {
// 		value, ok := map[string]ReportRange{
// 			"Today":      Today,
// 			"Yesterday":  Yesterday,
// 			"LastWeek":   LastWeek,
// 			"ThisWeek":   ThisWeek,
// 			"Custom":     Custom,
// 			"EYesterday": EYesterday,
// 		}[args[0]]
// 		if ok {
// 			return value
// 		}
// 	}
// 	return Yesterday
// }

// func Main() {
// 	// Sample args for testing
// 	args := []string{"Today"}
// 	reportRange := GetReportRange(args)

// 	startDateTime := time.Date(2023, 9, 14, 0, 0, 0, 0, time.UTC)
// 	endDateTime := startDateTime.AddDate(0, 0, 1)

// 	switch reportRange {
// 	case Today:
// 		startDateTime = time.Now().UTC().Truncate(24 * time.Hour)
// 		endDateTime = startDateTime
// 	case Yesterday:
// 		startDateTime = time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
// 		endDateTime = startDateTime
// 	case EYesterday:
// 		startDateTime = time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 0, -2)
// 		endDateTime = startDateTime
// 	case ThisWeek:
// 		today := time.Now().UTC()
// 		offset := int(time.Saturday) - int(today.Weekday())
// 		endDateTime = today.AddDate(0, 0, offset)
// 		startDateTime = endDateTime.AddDate(0, 0, -6)
// 	case LastWeek:
// 		today := time.Now().UTC()
// 		offset := int(time.Sunday) - int(today.Weekday()) - 1
// 		endDateTime = today.AddDate(0, 0, offset)
// 		startDateTime = endDateTime.AddDate(0, 0, -6)
// 	}

// 	timeEntries, err := FetchTimeEntries(startDateTime, endDateTime)
// 	if err != nil {
// 		fmt.Println("error fetching time entries:", err)
// 	}

// 	includeProjectAndClient(timeEntries)

// 	rpt := &Report{
// 		StartDate:       startDateTime,
// 		EndDate:         endDateTime,
// 		ProjectSections: make([]*ProjectSection, 0),
// 	}

// 	// Assuming ProjectKey, ParseTitle, and other related functions and types are defined
// 	// Group by ProjectKey
// 	projectGroups := make(map[ProjectKey]TogglTimeEntries)
// 	for _, entry := range timeEntries {
// 		key := ProjectKey{
// 			Client:    entry.Client,
// 			Project:   entry.Project,
// 			ProjectID: entry.ProjectID,
// 		}
// 		projectGroups[key] = append(projectGroups[key], entry)
// 	}

// 	// Sort and process each group
// 	for key, entries := range projectGroups {
// 		titleParseResult, err := ParseProjectTitle(key.Project)
// 		if err != nil {
// 			fmt.Println("Error parsing title:", err)
// 			continue
// 		}

// 		ps := &ProjectSection{
// 			Client:      key.Client,
// 			Project:     key.Project,
// 			ProjectID:   fmt.Sprintf("%d", key.ProjectID),
// 			IsTask:      titleParseResult.IsTask,
// 			TaskID:      titleParseResult.TaskID,
// 			DaySections: make([]*DaySection, 0),
// 		}

// 		rpt.ProjectSections = append(rpt.ProjectSections, ps)

// 		// Group by day
// 		dayGroups := make(map[time.Time]TogglTimeEntries)
// 		for _, entry := range entries {
// 			day := entry.Start.Truncate(24 * time.Hour)
// 			dayGroups[day] = append(dayGroups[day], entry)
// 		}

// 		// Sort and process each day group
// 		for day, dayEntries := range dayGroups {
// 			ds := &DaySection{
// 				Date:             day,
// 				IntervalSections: make([]*IntervalSection, 0),
// 			}
// 			ps.DaySections = append(ps.DaySections, ds)

// 			// Sort entries by start time
// 			sort.Slice(dayEntries, func(i, j int) bool {
// 				return dayEntries[i].Start.Before(dayEntries[j].Start)
// 			})

// 			for _, entry := range dayEntries {
// 				durationHours := float64(entry.Duration) / 3600.0 // Convert seconds to hours
// 				is := &IntervalSection{
// 					StartDateTime: entry.Start,
// 					EndDateTime:   entry.Stop,
// 					Description:   entry.Description,
// 					Duration:      durationHours,
// 				}
// 				ds.IntervalSections = append(ds.IntervalSections, is)
// 			}
// 		}
// 	}

// }
