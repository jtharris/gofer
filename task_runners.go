package main

import (
	"fmt"
	"sync"

	"github.com/fatih/color"
)

type serialTaskRunner struct {
	quiet bool
	task  *GoferTask
}

func (r serialTaskRunner) Run() GoferTaskResult {
	taskResult := GoferTaskResult{
		commands: make([]GoferCommandResult, len(r.task.Commands)),
	}

	for i, command := range r.task.Commands {
		if !r.quiet {
			fmt.Print(command, " ... ")
		}

		result := RunCommand(command)
		taskResult.commands[i] = result

		if !r.quiet {
			if result.err != nil {
				color.Red("\u2717")
			} else {
				color.Green("\u2713")
			}
		}

		if result.err != nil {
			break
		}
	}

	return taskResult
}

type parallelTaskRunner struct {
	quiet bool
	task  *GoferTask
}

func (r parallelTaskRunner) Run() GoferTaskResult {
	var wg sync.WaitGroup
	taskResult := GoferTaskResult{
		commands: make([]GoferCommandResult, len(r.task.Commands)),
	}

	run := func(slot int, command string) {
		result := RunCommand(command)
		taskResult.commands[slot] = result

		// TODO:  Is this the right place for this?
		if !r.quiet {
			if result.err != nil {
				fmt.Println(command, "...", color.RedString("\u2717"))
			} else {
				fmt.Println(command, "...", color.GreenString("\u2713"))
			}
		}

		wg.Done()
	}

	wg.Add(len(r.task.Commands))
	for i, command := range r.task.Commands {
		go run(i, command)
	}

	wg.Wait()

	return taskResult
}

type explainTaskRunner struct {
	task *GoferTask
}

func (r explainTaskRunner) Run() GoferTaskResult {
	taskColor := color.New(color.FgHiYellow, color.Bold)
	taskColor.Println(r.task.Name)
	var preCommandColor *color.Color
	var preCommandString string

	if r.task.Parallel {
		preCommandColor = color.New(color.FgGreen, color.Bold)
		preCommandString = "  | "
	} else {
		preCommandColor = color.New(color.FgWhite, color.Bold)
		preCommandString = "  |> "
	}

	commandColor := color.New(color.FgWhite)

	for _, command := range r.task.Commands {
		preCommandColor.Print(preCommandString)
		commandColor.Println(command)
	}

	return GoferTaskResult{}
}
