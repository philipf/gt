package client

import (
	"github.com/philipf/gt/cmd/toggl"
	"github.com/spf13/cobra"
)

// gtdCmd represents the action command
var clientsCmd = &cobra.Command{
	Use:     "client",
	Aliases: []string{"clients"},
	Short:   "Client CRUD",
	Long:    `Manages clients using the Toggl API`,

	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	toggl.TogglCmd.AddCommand(clientsCmd)
}
