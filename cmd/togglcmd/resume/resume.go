package resume

import (
	"github.com/philipf/gt/cmd/togglcmd"
	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/cobra"
)

var timeService toggl.TimeService

var stopCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume the current time entry",
	Long:  `Resume the current time entry`,

	Run: func(cmd *cobra.Command, args []string) {
		err := timeService.ResumeLast()
		cobra.CheckErr(err)
	},
}

func init() {
	cobra.OnInitialize(initServices)
	togglcmd.TogglCmd.AddCommand(stopCmd)

	// timeCmd.Flags().BoolP("today", "", false, "Run a report for today")
}

func initServices() {
	timeService = initialiseTimeService()
}
