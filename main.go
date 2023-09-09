package main

import (
	"fmt"
	"os"

	"github.com/philipf/gt/cmd"
)

func main() {
	// Main entry point for the GT CLI and calls the root command
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
