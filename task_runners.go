package main

import (
	"log"
	"sync"
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
		result := RunCommand(command)
		taskResult.commands[i] = result

		if !r.quiet {
			log.Println(result.output)
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
			log.Println(result.output)
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
	log.Println(r.task.Name)
	for _, command := range r.task.Commands {
		log.Println("  |>", command)
	}

	return GoferTaskResult{}
}
