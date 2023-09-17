package project

import (
	"fmt"

	"github.com/philipf/gt/internal/togglservices"
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

		items, err := togglservices.GetProjects(nil)
		cobra.CheckErr(err)

		for _, i := range items {
			fmt.Printf("%d - %s\n", i.ID, i.Name)
		}
	},
}

func init() {
	projectsCmd.AddCommand(listCmd)

	// filter
	//listCmd.Flags().StringP("filter", "f", "", "Filter clients by name")

}
