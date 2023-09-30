package report

import (
	"fmt"
	"os"

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
		sd, ed := togglcmd.GetDateRange(cmd)

		rpt, err := reportService.BuildReport(sd, ed)
		cobra.CheckErr(err)

		textReport, err := (cmd.Flags().GetBool("text"))
		cobra.CheckErr(err)

		jsonReport, err := (cmd.Flags().GetBool("json"))
		cobra.CheckErr(err)

		if textReport || textOutput != "" || (!textReport && !jsonReport) {
			r, err := reportService.BuildStringReport(rpt)
			cobra.CheckErr(err)

			if textReport {
				fmt.Println(*r)
			}

			if textOutput != "" {
				os.WriteFile(textOutput, []byte(*r), 0644)
			}
		}

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

	reportCmd.Flags().BoolP("text", "t", false, "Produces a text report")
	reportCmd.Flags().BoolP("json", "j", false, "Produces a JSON report")

	reportCmd.Flags().StringVarP(&textOutput, "ot", "", "", "Writes a text report file.")
	reportCmd.Flags().Lookup("ot").NoOptDefVal = "/tmp/time.txt"

	reportCmd.Flags().StringVarP(&jsonOutput, "oj", "", "", "Writes a JSON report file.")
	reportCmd.Flags().Lookup("oj").NoOptDefVal = "/tmp/time.json"

	reportCmd.Flags().BoolP("today", "", false, "Run a report for today")
	reportCmd.Flags().BoolP("yesterday", "", false, "Run a report for yesterday")
	reportCmd.Flags().BoolP("eyesterday", "", false, "Run a report for eyesterday")
	reportCmd.Flags().BoolP("thisweek", "", false, "Run a report for this week")
	reportCmd.Flags().BoolP("lastweek", "", false, "Run a report for last week")

	reportCmd.Flags().StringP("startDate", "s", "", "Start date")
	reportCmd.Flags().StringP("endDate", "e", "", "End date")

}

func initServices() {
	reportService = initialiseReportService()
}
