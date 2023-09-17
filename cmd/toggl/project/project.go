package project

import (
	"github.com/philipf/gt/cmd/toggl"
	"github.com/spf13/cobra"
)

// gtdCmd represents the action command
var projectsCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"projects"},
	Short:   "Project CRUD",
	Long:    `Manages projects using the Toggl API`,

	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	toggl.TogglCmd.AddCommand(projectsCmd)
}
