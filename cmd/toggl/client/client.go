package client

import (
	cmdToggl "github.com/philipf/gt/cmd/toggl"
	internalToggl "github.com/philipf/gt/internal/toggl"
	"github.com/spf13/cobra"
)

var clientService internalToggl.ClientService

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
	cobra.OnInitialize(initServices)
	cmdToggl.TogglCmd.AddCommand(clientsCmd)
}

func initServices() {
	clientService = initializeClientService()
}
