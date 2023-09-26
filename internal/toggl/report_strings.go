package toggl

import (
	"fmt"
	"strings"
)

type StringReport struct {
}

func NewStringReport() StringReport {
	return StringReport{}
}

func (*StringReport) Generate(r *Report) *string {
	totalDuration := 0.0
	var sb strings.Builder

	sd := r.StartDate.Local()
	ed := r.EndDate.Local()

	// Print the header with the date range
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("Start: %s to %s\n", sd.Format("2006-01-02"), ed.Format("2006-01-02")))
	sb.WriteString("-------------------------------\n")

	for _, projectSection := range r.ProjectSections {
		sb.WriteString(fmt.Sprintf("%s\n", projectSection.Project))

		for _, daySection := range projectSection.DaySections {
			// Print DaySection
			d := daySection.Date
			sb.WriteString(fmt.Sprintf("%s %s\n", d.Format("2006-01-02"), d.Weekday().String()))
			sb.WriteString("Start   End    Time   Notes\n")

			for _, intervalSection := range daySection.IntervalSections {
				// Print IntervalSection
				sb.WriteString(fmt.Sprintf("%s - %s  %s  %s\n",
					intervalSection.StartDateTime.Format("15:04"),
					intervalSection.EndDateTime.Format("15:04"),
					intervalSection.DurationS(),
					intervalSection.Description,
				))
			}
			dayDuration := roundToQuarterHour(daySection.Duration())
			totalDuration += dayDuration
			sb.WriteString(fmt.Sprintf("Duration: %.2f\n\n", dayDuration))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("Report Duration: %.2f\n", totalDuration))
	//sb.WriteString(fmt.Sprintf("Report Duration: %.2f\n", r.Duration()))
	sb.WriteString("---EOF---")

	str := sb.String()

	return &str
}
