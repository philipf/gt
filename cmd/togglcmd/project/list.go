package project

import (
	"fmt"

	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	Long: `Lists projects using the Toggl API

Example:
gt toggl project list
`,

	Run: func(cmd *cobra.Command, args []string) {
		includeArchived, err := cmd.Flags().GetBool("includeArchived")
		cobra.CheckErr(err)

		clientId, err := cmd.Flags().GetInt64("clientId")
		cobra.CheckErr(err)

		name, err := cmd.Flags().GetString("name")
		cobra.CheckErr(err)

		filter := toggl.GetProjectsOpts{
			IncludeArchived: includeArchived,
			ClientID:        clientId,
			Name:            name,
		}

		items, err := projectService.Get(&filter)
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
			if !projectService.HasValidName(i.Name) {
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
	projectCmd.AddCommand(listCmd)

	// filter
	listCmd.Flags().Bool("validate", false, "Validate projects matches the naming convention")
	listCmd.Flags().Bool("includeArchived", false, "Include archived projects")
	listCmd.Flags().Int64P("clientId", "c", 0, "Filter by client ID")
	listCmd.Flags().StringP("name", "n", "", "Filter by project name")
}
