/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// settingsCmd represents the settings command
var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Prints out the settings",
	Long: `Settings are stored in a config.yaml file in the user's home directory under the .gt directory.
	If the file does not exist, it will be created with default values.
	It is possible to override the default values by creating a config.yaml file in the current directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		// print out the current settings
		fmt.Println("Current settings:")
		for _, key := range viper.AllKeys() {
			fmt.Printf("%s: %s\n", key, viper.Get(key))
		}
	},
}

func init() {
	RootCmd.AddCommand(settingsCmd)
}
