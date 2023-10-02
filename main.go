package main

import (
	"fmt"
	"os"

	"github.com/philipf/gt/cmd"

	_ "github.com/philipf/gt/cmd/cal"
	_ "github.com/philipf/gt/cmd/gtd"
	_ "github.com/philipf/gt/cmd/togglcmd"
	_ "github.com/philipf/gt/cmd/togglcmd/client"
	_ "github.com/philipf/gt/cmd/togglcmd/project"
	_ "github.com/philipf/gt/cmd/togglcmd/report"
	_ "github.com/philipf/gt/cmd/togglcmd/resume"
	_ "github.com/philipf/gt/cmd/togglcmd/stop"
	_ "github.com/philipf/gt/cmd/togglcmd/time"
)

func main() {
	// Main entry point for the GT CLI and calls the root command
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
