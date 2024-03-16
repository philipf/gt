package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of the gt cli",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("gt version 0.8.1")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
