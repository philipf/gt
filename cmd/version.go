package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of the gt cli",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("gt version 0.10.4")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
