package cmd

import (
	"fmt"
	"os"

	"github.com/philipf/gt/internal/settings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "gt - Go Time by Philip Fourie",
	Long:  `Provides multiple utilties for enhancing daily tasks`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmd.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	err := settings.Init()
	cobra.CheckErr(err)

	fmt.Println("setting open AI key")
	os.Setenv("OPENAI_API_KEY", viper.GetString("ai.openAiKey"))

}
