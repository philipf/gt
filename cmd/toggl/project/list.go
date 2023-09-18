package project

import (
	"fmt"

	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/cobra"
)

// gtdCmd represents the action command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	Long: `Lists projects using the Toggl API

Example:
gt toggl project list
`,

	Run: func(cmd *cobra.Command, args []string) {
		//filter, _ := cmd.Flags().GetString("filter")

		service := toggl.ProjectServiceImplementation{}
		items, err := service.GetProjects(nil)
		cobra.CheckErr(err)

		validate, err := cmd.Flags().GetBool("validate")
		cobra.CheckErr(err)

		if !validate {
			for _, i := range items {
				fmt.Printf("%d - %s\n", i.ID, i.Name)
			}
			return
		}

		invalidProjects := make(toggl.TogglProjects, 0)

		for _, i := range items {
			// only print the projects that don't match the naming convention
			if !service.ValidProjectName(i.Name) {
				invalidProjects = append(invalidProjects, i)
			}
		}

		if len(invalidProjects) == 0 {
			fmt.Printf("All %d projects match the naming convention\n", len(items))
			return
		} else {
			fmt.Printf("%d project(s) do not match the naming convention:\n", len(invalidProjects))

			for _, i := range invalidProjects {
				fmt.Printf("%d - %s\n", i.ID, i.Name)
			}
		}

	},
}

func init() {
	projectsCmd.AddCommand(listCmd)

	// filter
	listCmd.Flags().Bool("validate", false, "Validate projects matches the naming convention")
}
