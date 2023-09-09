package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of the GT CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GT version 0.0.10")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}