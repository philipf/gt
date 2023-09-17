package toggl

import (
	"github.com/philipf/gt/cmd"
	"github.com/spf13/cobra"
)

// gtdCmd represents the action command
var TogglCmd = &cobra.Command{
	Use:     "toggl",
	Short:   "Toggl time tracking commands",
	Long:    `Toggl commands for time tracking, project and client creation`,
	Aliases: []string{"toggle"},

	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	cmd.RootCmd.AddCommand(TogglCmd)
}
