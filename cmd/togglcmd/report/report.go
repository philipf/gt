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

		rpt, err := reportService.BuildReport(sd, ed)
		if err != nil {
			cobra.CheckErr(err)
		}

		console, err := (cmd.Flags().GetBool("console"))
		if err != nil {
			cobra.CheckErr(err)
		}

		if console {
			r, err := reportService.BuildStringReport(rpt)
			if err != nil {
				cobra.CheckErr(err)
			}
			fmt.Println(*r)
		}

		json, err := (cmd.Flags().GetBool("json"))
		if err != nil {
			cobra.CheckErr(err)
		}

		if json {
			jsonOutput, err := cmd.Flags().GetString("jsonOutput")
			if err != nil {
				cobra.CheckErr(err)
			}

			jsonBytes, err := reportService.BuildJsonReport(rpt)
			if err != nil {
				cobra.CheckErr(err)
			}

			os.WriteFile(jsonOutput, *jsonBytes, 0644)
		}
	},
}

func init() {
	cobra.OnInitialize(initServices)
	togglcmd.TogglCmd.AddCommand(reportCmd)

	reportCmd.Flags().StringP("sd", "s", "", "Start date")
	reportCmd.Flags().StringP("ed", "e", "", "End date")

	reportCmd.Flags().BoolP("console", "c", false, "Console Report")
	reportCmd.Flags().BoolP("json", "j", false, "JSON Report")

	reportCmd.Flags().StringP("jsonOutput", "o", "/tmp/time.json", "JSON output file")

	reportCmd.MarkFlagRequired("sd")
	reportCmd.MarkFlagRequired("ed")

}

func initServices() {
	reportService = initialiseReportService()
}
