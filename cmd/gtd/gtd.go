package gtd

import (
	"github.com/philipf/gt/cmd"
	"github.com/spf13/cobra"
)

var UseAi bool

// gtdCmd represents the action command
var gtdCmd = &cobra.Command{
	Use:   "gtd",
	Short: "Create a new GTD action",
	Long: `Create a new GTD action. This will create a new action in the inbox and add a new todo to the kanban board.
If no description is provided, only the todo will be added to the kanban board.
Multi line input is supported for the description. To end the description, enter a full stop on a new line.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		addCmd.Run(cmd, args)
	},
}

func init() {
	cmd.RootCmd.AddCommand(gtdCmd)

	gtdCmd.PersistentFlags().BoolVarP(&UseAi, "ai", "", false, "Use AI to assist with the creation of a new action")
}
