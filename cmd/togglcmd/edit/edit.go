package edit

import (
	"fmt"

	"github.com/philipf/gt/cmd/togglcmd"
	"github.com/philipf/gt/internal/console"
	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/cobra"
)

var timeService toggl.TimeService

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the description of the running time entry",
	Long:  `Edit the description of the running time entry`,

	Run: func(cmd *cobra.Command, _ []string) {
		// If the description flag is set, then update the description of the running time entry
		desc, err := cmd.Flags().GetString("description")
		cobra.CheckErr(err)

		if desc == "" {
			desc, err = getUserInput("Description:")
			cobra.CheckErr(err)
		}

		err = timeService.EditDesc(desc)
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
	cobra.OnInitialize(initServices)
	togglcmd.TogglCmd.AddCommand(editCmd)

	editCmd.Flags().StringP("description", "d", "", "Description of the time entry")
}

func initServices() {
	timeService = initialiseTimeService()
}

func getUserInput(prompt string) (string, error) {
	fmt.Println(prompt)
	return console.ReadLine()
}
