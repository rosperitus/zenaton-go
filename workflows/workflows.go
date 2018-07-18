package workflows

import (
	"fmt"

	"github.com/zenaton/zenaton-go/tasks"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/version"
	"github.com/zenaton/zenaton-go/v1/zenaton/wait"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var (
	SequentialWorkflow = workflow.Workflow{
		Name: "SequentialWorkflow",
		// this is ugly to have it have a return value
		HandleFunc: func() interface{} {
			tasks.TaskA.Execute()
			tasks.TaskB.Execute()
			return nil
		},
	}

	AsynchronousWorkflow = workflow.Workflow{
		Name: "AsynchronousWorkflow",
		HandleFunc: func() interface{} {
			tasks.TaskA.Dispatch()
			tasks.TaskB.Execute()
			return nil
		},
	}

	ParallelWorkflow = workflow.Workflow{
		Name: "ParallelWorkflow",
		HandleFunc: func() interface{} {
			runParallel := task.Tasks{
				tasks.TaskA,
				tasks.TaskB, tasks.TaskB, tasks.TaskB, tasks.TaskB,
			}
			outcomes := runParallel.Execute()
			fmt.Println("outcomes: ", outcomes)
			tasks.TaskC.Execute()
			return nil
		},
	}

	EventWorkflow = workflow.Workflow{
		Name: "EventWorkflow",
		HandleFunc: func() interface{} {
			tasks.TaskA.Execute()
			tasks.TaskB.Execute()
			//todo: ugly to have to return nil all the time, can I do better?
			return nil
		},
		OnEvent: func(eventName string, eventData interface{}) {
			if eventName == "MyEvent" {
				tasks.TaskC.Execute()
			}
		},
		//todo: do something sensible when you don't have an ID function
		ID: func() string {
			return "MyId"
		},
	}

	WaitWorkflow = workflow.Workflow{
		Name: "WaitWorkflow",
		HandleFunc: func() interface{} {
			// todo: figure out how to do something like this.email in javascript example
			tasks.TaskA.Execute()
			// todo: kind of ugly to pass in nil here
			wait.New(nil).Seconds(5).Execute()
			tasks.TaskB.Execute()
			return nil
		},
	}

	WaitEventWorkflow = workflow.Workflow{
		Name: "WaitEventWorkflow",
		HandleFunc: func() interface{} {

			// Wait until the event or 4 seconds
			event := wait.New("MyEvent").Seconds(4).Execute()

			// If event has been triggered
			if event != nil {
				// Execute TaskB
				tasks.TaskA.Execute()
			} else {
				// Execute Task B
				tasks.TaskB.Execute()
			}
			return nil
		},

		ID: func() string {
			return "MyId"
		},
	}

	VersionWorkflow = version.New("VersionWorkflow", []*workflow.Workflow{
		VersionWorkflow_v0,
		VersionWorkflow_v1,
		VersionWorkflow_v2,
	})

	VersionWorkflow_v0 = &workflow.Workflow{
		Name: "VersionWorkflow_v0",
		HandleFunc: func() interface{} {
			task.Tasks{
				tasks.TaskA,
				tasks.TaskB,
			}.Execute()
			return nil
		},
	}

	VersionWorkflow_v1 = &workflow.Workflow{
		Name: "VersionWorkflow_v1",
		HandleFunc: func() interface{} {
			task.Tasks{
				tasks.TaskA,
				tasks.TaskB,
				tasks.TaskC,
			}.Execute()
			return nil
		},
	}

	VersionWorkflow_v2 = &workflow.Workflow{
		Name: "VersionWorkflow_v2",
		HandleFunc: func() interface{} {
			task.Tasks{
				tasks.TaskA,
				tasks.TaskB,
				tasks.TaskC,
				tasks.TaskD,
			}.Execute()
			return nil
		},
	}
)
