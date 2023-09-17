package main

import (
	"fmt"
	"os"

	"github.com/philipf/gt/cmd"

	_ "github.com/philipf/gt/cmd/cal"
	_ "github.com/philipf/gt/cmd/gtd"
	_ "github.com/philipf/gt/cmd/toggl"
	_ "github.com/philipf/gt/cmd/toggl/client"
	_ "github.com/philipf/gt/cmd/toggl/project"
)

func main() {
	// Main entry point for the GT CLI and calls the root command
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
