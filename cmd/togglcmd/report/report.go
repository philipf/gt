package report

import (
	"fmt"
	"os"
	"time"

	"github.com/philipf/gt/cmd/togglcmd"
	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/cobra"
)

var reportService toggl.ReportService
var textOutput string
var jsonOutput string

// gtdCmd represents the action command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Report Service",
	Long:  `Report Service	`,

	Run: func(cmd *cobra.Command, args []string) {
		sd, ed := getDateRange(cmd)

		rpt, err := reportService.BuildReport(sd, ed)
		cobra.CheckErr(err)

		textReport, err := (cmd.Flags().GetBool("text"))
		cobra.CheckErr(err)

		if textReport || textOutput != "" {
			r, err := reportService.BuildStringReport(rpt)
			cobra.CheckErr(err)

			if textReport {
				fmt.Println(*r)
			}

			if textOutput != "" {
				os.WriteFile(textOutput, []byte(*r), 0644)
			}
		}

		jsonReport, err := (cmd.Flags().GetBool("json"))
		cobra.CheckErr(err)

		if jsonReport || jsonOutput != "" {
			jsonBytes, err := reportService.BuildJsonReport(rpt)
			cobra.CheckErr(err)

			if jsonReport {
				fmt.Println(string(*jsonBytes))
			}

			if jsonOutput != "" {
				os.WriteFile(jsonOutput, *jsonBytes, 0644)
			}
		}
	},
}

func init() {
	cobra.OnInitialize(initServices)
	togglcmd.TogglCmd.AddCommand(reportCmd)

	reportCmd.Flags().StringP("startDate", "s", "", "Start date")
	reportCmd.Flags().StringP("endDate", "e", "", "End date")

	reportCmd.Flags().BoolP("text", "t", false, "Produces a text report")
	reportCmd.Flags().BoolP("json", "j", false, "Produces a JSON report")

	reportCmd.Flags().StringVarP(&textOutput, "ot", "", "", "Writes a text report file.")
	reportCmd.Flags().Lookup("ot").NoOptDefVal = "/tmp/time.txt"

	reportCmd.Flags().StringVarP(&jsonOutput, "oj", "", "", "Writes a JSON report file.")

	reportCmd.Flags().BoolP("today", "", false, "Run a report for today")
	reportCmd.Flags().BoolP("eyesterday", "", false, "Run a report for eyesterday")
	reportCmd.Flags().BoolP("yesterday", "", false, "Run a report for yesterday")
	reportCmd.Flags().BoolP("thisweek", "", false, "Run a report for this week")
	reportCmd.Flags().BoolP("lastweek", "", false, "Run a report for last week")

	//	reportCmd.MarkFlagRequired("sd")
	//	reportCmd.MarkFlagRequired("ed")

}

func initServices() {
	reportService = initialiseReportService()
}

func getDateRange(cmd *cobra.Command) (time.Time, time.Time) {
	// First check if any of the fixed periods have been selected (today, yesterday, etc)
	today, err := cmd.Flags().GetBool("today")
	if err != nil {
		cobra.CheckErr(err)
	}

	now := time.Now()

	if today {
		// Set the start date to today, trimming the time
		sd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		// Set the end date to tomorrow, trimming the time and subtracting a second
		ed := sd.AddDate(0, 0, 1).Add(-time.Second)
		return sd, ed
	}

	// check if yesterday has been selected
	yesterday, err := cmd.Flags().GetBool("yesterday")
	if err != nil {
		cobra.CheckErr(err)
	}

	if yesterday {
		// Set the start date to yesterday, trimming the time
		now := now.AddDate(0, 0, -1)
		sd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		// Set the end date to today, trimming the time and subtracting a second
		ed := sd.AddDate(0, 0, 1).Add(-time.Second)
		return sd, ed
	}

	// check if eyesterday has been selected
	eyesterday, err := cmd.Flags().GetBool("eyesterday")
	if err != nil {
		cobra.CheckErr(err)
	}

	if eyesterday {
		// Set the start date to yesterday, trimming the time
		now := now.AddDate(0, 0, -2)
		sd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		// Set the end date to yesterday, trimming the time and subtracting a second
		ed := sd.AddDate(0, 0, 1).Add(-time.Second)
		return sd, ed
	}

	// check if thisweek has been selected
	thisweek, err := cmd.Flags().GetBool("thisweek")
	if err != nil {
		cobra.CheckErr(err)
	}

	if thisweek {
		// TODO: get logic from the gt-at package
		thisSunday := sundayOfTheWeek(time.Now())
		sd := time.Date(thisSunday.Year(), thisSunday.Month(), thisSunday.Day(), 0, 0, 0, 0, thisSunday.Location())
		ed := sd.AddDate(0, 0, 7).Add(-time.Second)

		return sd, ed
	}

	// check if lastweek has been selected
	lastweek, err := cmd.Flags().GetBool("lastweek")
	if err != nil {
		cobra.CheckErr(err)
	}

	if lastweek {
		// TODO: get logic from the gt-at package
		thisSunday := sundayOfTheWeek(time.Now())
		lastSunday := thisSunday.AddDate(0, 0, -7)
		sd := time.Date(lastSunday.Year(), lastSunday.Month(), lastSunday.Day(), 0, 0, 0, 0, lastSunday.Location())
		ed := sd.AddDate(0, 0, 7).Add(-time.Second)

		return sd, ed
	}

	sdStr, err := cmd.Flags().GetString("startDate")
	if err != nil {
		cobra.CheckErr(err)
	}

	edStr, err := cmd.Flags().GetString("endDate")
	if err != nil {
		cobra.CheckErr(err)
	}

	loc, err := time.LoadLocation("Local")
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

	return sd, ed
}

// SundayOfTheWeek returns the date of the Sunday of the week based on a provided date.
func sundayOfTheWeek(t time.Time) time.Time {
	// Subtract the weekday number from the given date.
	// Since Sunday = 0 in time.Weekday, it gives the exact offset we need.
	offset := int(t.Weekday())
	return t.AddDate(0, 0, -offset)
}
