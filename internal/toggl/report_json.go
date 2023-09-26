package toggl

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

type GtAutoTaskRequest struct {
	ID        int       `json:"id"`
	IsTicket  bool      `json:"isTicket"`
	Date      time.Time `json:"date"`
	StartTime string    `json:"startTime"`
	Duration  float64   `json:"duration"`
	Summary   string    `json:"summary"`
	Project   string    `json:"project"`
}

type JsonReport struct {
}

func NewJsonReport() JsonReport {
	return JsonReport{}
}

func (*JsonReport) Generate(report *Report) (*[]byte, error) {
	var result []GtAutoTaskRequest

	for _, ps := range report.ProjectSections {
		for _, ds := range ps.DaySections {
			intervalSections := ds.IntervalSections
			sort.Slice(intervalSections, func(i, j int) bool {
				return intervalSections[i].StartDateTime.Before(intervalSections[j].StartDateTime)
			})

			var sb strings.Builder

			sb.WriteString("Start   End    Time   Notes\n")
			for _, ins := range intervalSections {
				sb.WriteString(fmt.Sprintf("%s - %s  %s  %s\n",
					ins.StartDateTime.Format("15:04"),
					ins.EndDateTime.Format("15:04"),
					ins.DurationS(),
					ins.Description))
			}

			duration := roundToQuarterHour(ds.Duration())
			sb.WriteString(fmt.Sprintf("Duration: %.2f", duration))

			r := GtAutoTaskRequest{
				ID:        ps.TaskID,
				IsTicket:  ps.IsTask,
				Date:      ds.Date,
				StartTime: intervalSections[0].StartDateTime.Format("15:04"),
				Duration:  duration,
				Summary:   sb.String(),
				Project:   ps.Project,
			}

			result = append(result, r)
		}
	}

	bytes, err := json.MarshalIndent(result, "", "    ")

	if err != nil {
		return nil, err
	}

	return &bytes, nil
}
