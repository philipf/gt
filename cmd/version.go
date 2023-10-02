package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of the gt cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gt version 0.4.0")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
