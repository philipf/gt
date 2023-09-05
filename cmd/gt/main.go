// The main entry point for the gt cli command
package main

import (
	"fmt"
	"time"

	alltasks "github.com/philipf/gt/cmd/alltasks"
	"github.com/philipf/gt/internal/gtd"
	"github.com/philipf/gt/internal/tasks"
)

func main() {
	fmt.Println("GT Version 0.0.4")

	dt := time.Now()
	externalId := ""
	externalLink := ""

	title := "get cli"
	description := "create a cli for gt"

	action, err := gtd.CreateAction(externalId, title, description, externalLink, dt, dt, gtd.In)
	if err != nil {
		panic(err)
	}

	fmt.Println(action)

	err = gtd.ActionToMd(action, "./inbox/")
	if err != nil {
		panic(err)
	}
}

func OldMain() {
	tsks := alltasks.GetTasks()

	taskList := []tasks.Task{}

	for _, t := range tsks {
		var externalLink string = ""

		task, err := tasks.CreateTask(*t.GetId(), *t.GetTitle(), *t.GetBody().GetContent(), externalLink, *t.GetCreatedDateTime(), *t.GetLastModifiedDateTime())

		//task.DueAt = *t.GetDueDateTime().

		if err != nil {
			panic(err)
		}
		taskList = append(taskList, *task)

	}

	// list of actions
	//actions := []gtd.Action{}

	actions, err := gtd.MapTasks(taskList)
	if err != nil {
		panic(err)
	}

	// for i := 0; i < 5; i++ {
	// 	title := fmt.Sprintf("Test Action %d", i)
	// 	action, err := gtd.CreateAction("ex", title, "Test Description", "Test Link", time.Now(), time.Now(), gtd.In)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	actions = append(actions, *action)

	// 	fmt.Println(action)
	// }

	gtd.ActionsToMd(actions, "./inbox/")
}
