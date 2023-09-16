package cmd

import (
	"github.com/philipf/gt/internal/settings"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "gt - Go Time by Philip Fourie",
	Long:  `gt - Got Time, provides utilties for enhancing daily tasks`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	err := settings.Init()
	cobra.CheckErr(err)
}
