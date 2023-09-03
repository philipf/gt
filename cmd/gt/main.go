// The main entry point for the gt cli command
package main

import (
	"fmt"
	"time"

	"github.com/philipf/gt/internal/gtd"
)

func main() {
	fmt.Println("GT Version 0.0.3")

	// Create dummy actions for testing

	// list of actions
	actions := []gtd.Action{}

	for i := 0; i < 5; i++ {
		title := fmt.Sprintf("Test Action %d", i)
		action, err := gtd.CreateAction("ex", title, "Test Description", "Test Link", time.Now(), time.Now(), gtd.In)
		if err != nil {
			panic(err)
		}

		actions = append(actions, *action)

		fmt.Println(action)
	}

	gtd.ExportToMd(actions, "./inbox/")
}
