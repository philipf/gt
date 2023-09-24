package report

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/philipf/gt/cmd/togglcmd"
	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/cobra"
)

var reportService toggl.ReportService

// gtdCmd represents the action command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Report Service",
	Long:  `Report Service	`,

	Run: func(cmd *cobra.Command, args []string) {
		loc, err := time.LoadLocation("Local")
		if err != nil {
			fmt.Println("Error loading location:", err)
			return
		}

		sdStr, err := cmd.Flags().GetString("sd")
		if err != nil {
			cobra.CheckErr(err)
		}

		edStr, err := cmd.Flags().GetString("ed")
		if err != nil {
			cobra.CheckErr(err)
		}

		sd, err := time.ParseInLocation("2006/01/02", sdStr, loc)
		if err != nil {
			cobra.CheckErr(err)
		}

		ed, err := time.ParseInLocation("2006/01/02", edStr, loc)
		if err != nil {
			cobra.CheckErr(err)
		}

		ed = ed.AddDate(0, 0, 1).Add(-time.Second)

		fmt.Println(sd)
		fmt.Println(ed)

		r, err := reportService.GetReport(sd, ed)
		if err != nil {
			cobra.CheckErr(err)
		}

		fmt.Println(generateReport(r))
	},
}

func init() {
	cobra.OnInitialize(initServices)
	togglcmd.TogglCmd.AddCommand(reportCmd)

	reportCmd.Flags().StringP("sd", "s", "", "Start date")
	reportCmd.Flags().StringP("ed", "e", "", "End date")
}

func initServices() {
	reportService = initialiseReportService()
}

func generateReport(r *toggl.Report) string {
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
	sb.WriteString("---")

	return sb.String()
}

func roundToQuarterHour(duration float64) float64 {
	rounded := math.Ceil(duration*4) / 4.0
	return rounded
}
