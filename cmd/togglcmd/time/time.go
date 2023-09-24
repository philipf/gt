package time

import (
	"fmt"
	"time"

	"github.com/philipf/gt/cmd/togglcmd"
	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/cobra"
)

var timeService toggl.TimeService

// gtdCmd represents the action command
var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "Time entries",
	Long:  `Time entries`,

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

		entries, err := timeService.GetTimeEntries(sd, ed)

		if err != nil {
			cobra.CheckErr(err)
		}

		for _, i := range entries {
			fmt.Printf("%d: %v - [%s]:[%s]%s\n", i.ID, i.Start, i.Project, i.Client, i.Description)
		}
	},
}

func init() {
	cobra.OnInitialize(initServices)
	togglcmd.TogglCmd.AddCommand(timeCmd)

	timeCmd.Flags().StringP("sd", "s", "", "Start date")
	timeCmd.Flags().StringP("ed", "e", "", "End date")
}

func initServices() {
	timeService = initialiseTimeService()
}
