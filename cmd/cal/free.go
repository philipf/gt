package cal

import (
	"fmt"

	"github.com/philipf/gt/cmd"
	"github.com/spf13/cobra"
)

// freeCmd represents the free command
var freeCmd = &cobra.Command{
	Use:   "free",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("free called")
	},
}

func init() {
	cmd.RootCmd.AddCommand(freeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// freeCmd.PersistentFlags().String("foo", "", "A help for foo")

}
