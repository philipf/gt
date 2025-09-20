package time

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/philipf/gt/cmd/togglcmd"
	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/cobra"

	"github.com/olekukonko/tablewriter"
)

var timeService toggl.TimeService

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "Time entries",
	Long:  `Time entries`,

	Run: func(cmd *cobra.Command, args []string) {
		sd, ed := togglcmd.GetDateRange(cmd)

		entries, err := timeService.Get(sd, ed)
		cobra.CheckErr(err)

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Start.Before(entries[j].Start)
		})

		csv, err := cmd.Flags().GetBool("csv")
		cobra.CheckErr(err)

		if csv {
			generateCsv(entries)
			return
		}

		generateTable(entries)

	},
}

func generateTable(entries toggl.TogglTimeEntries) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "Start", "Stop", "Duration", "Project", "Description"})

	for _, entry := range entries {
		row := []string{
			fmt.Sprintf("%d", entry.ID),
			entry.Start.Format(time.DateTime),
			entry.Stop.Format(time.DateTime),
			fmt.Sprintf("%d", entry.Duration/60),
			entry.Project,
			entry.Description,
		}
		table.Append(row)
	}

	//table.AutoWrapText(false)
	table.Render()
}

func generateCsv(entries toggl.TogglTimeEntries) {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Write the CSV header
	writer.Write([]string{"ID", "Start", "Stop", "Duration", "Description", "Client", "Project"})

	// Write the CSV data
	for _, entry := range entries {
		row := []string{
			fmt.Sprintf("%d", entry.ID),
			entry.Start.Format(time.RFC3339),
			entry.Stop.Format(time.RFC3339),
			fmt.Sprintf("%d", entry.Duration),
			entry.Description,
			entry.Client,
			entry.Project,
		}
		err := writer.Write(row)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
}

func init() {
	cobra.OnInitialize(initServices)
	togglcmd.TogglCmd.AddCommand(timeCmd)

	timeCmd.Flags().BoolP("today", "", false, "Run a report for today")
	timeCmd.Flags().BoolP("yesterday", "", false, "Run a report for yesterday")
	timeCmd.Flags().BoolP("eyesterday", "", false, "Run a report for eyesterday")
	timeCmd.Flags().BoolP("thisweek", "", false, "Run a report for this week")
	timeCmd.Flags().BoolP("lastweek", "", false, "Run a report for last week")

	timeCmd.Flags().StringP("startDate", "s", "", "Start date")
	timeCmd.Flags().StringP("endDate", "e", "", "End date")

	timeCmd.Flags().BoolP("csv", "", false, "Export to CSV")
}

func initServices() {
	timeService = initialiseTimeService()
}
