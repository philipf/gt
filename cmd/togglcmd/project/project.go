package project

import (
	"github.com/philipf/gt/cmd/togglcmd"
	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/cobra"
)

var projectService toggl.ProjectService

// gtdCmd represents the action command
var projectCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"projects"},
	Short:   "Project CRUD",
	Long:    `Manages projects using the Toggl API`,

	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	cobra.OnInitialize(initServices)
	togglcmd.TogglCmd.AddCommand(projectCmd)
}

func initServices() {
	projectService = initialiseProjectService()
}
