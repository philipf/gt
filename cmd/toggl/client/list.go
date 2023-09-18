package client

import (
	"fmt"

	"github.com/philipf/gt/internal/toggl"
	"github.com/spf13/cobra"
)

// gtdCmd represents the action command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List clietns",
	Long: `Lists clients using the Toggl API

Example:
gt toggl client list
`,

	Run: func(cmd *cobra.Command, args []string) {
		filter, _ := cmd.Flags().GetString("filter")

		service := toggl.ClientServiceImplementation{}
		//service.ClientGateway = http.TogglClientGateway{}

		items, err := service.GetClients(filter)
		cobra.CheckErr(err)

		for _, i := range items {
			fmt.Printf("%d - %s\n", i.ID, i.Name)
		}
	},
}

func init() {
	clientsCmd.AddCommand(listCmd)

	// filter
	listCmd.Flags().StringP("filter", "f", "", "Filter clients by name")

}
