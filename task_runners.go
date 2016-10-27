package main

import "sync"

type serialTaskRunner struct {
	quiet bool
	task  *GoferTask
}

func (r serialTaskRunner) Run() GoferTaskResult {
	reporter := CommandReporter{
		quiet:    r.quiet,
		parallel: false,
	}
	taskResult := GoferTaskResult{
		commands: make([]GoferCommandResult, len(r.task.Commands)),
	}

	reporter.ReportTaskStart(r.task)

	for i, command := range r.task.Commands {
		result := RunCommand(command)
		taskResult.commands[i] = result
		reporter.ReportResult(&result)

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
	reporter := CommandReporter{
		quiet:    r.quiet,
		parallel: true,
	}
	taskResult := GoferTaskResult{
		commands: make([]GoferCommandResult, len(r.task.Commands)),
	}

	reporter.ReportTaskStart(r.task)

	run := func(slot int, command string) {
		result := RunCommand(command)
		taskResult.commands[slot] = result
		reporter.ReportResult(&result)

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
	reporter := CommandReporter{
		quiet:    false,
		parallel: r.task.Parallel,
	}

	reporter.ReportTaskStart(r.task)

	for _, command := range r.task.Commands {
		result := &GoferCommandResult{
			command: command,
		}
		reporter.ReportResult(result)
	}

	return GoferTaskResult{}
}
