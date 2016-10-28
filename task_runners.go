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
		log, err := r.task.CreateLogFile(i)
		if err != nil {
			taskResult.commands[i] = GoferCommandResult{
				command: command,
				err:     err,
			}

			break
		}

		if log != nil {
			defer log.Close()
		}

		result := RunCommand(command, log)
		taskResult.commands[i] = *result
		reporter.ReportResult(result)

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
		defer wg.Done()
		var result *GoferCommandResult

		// TODO:  DRY this up w/ the serialTaskRunner - lots of copy-paste
		log, err := r.task.CreateLogFile(slot)
		if err != nil {
			result = &GoferCommandResult{
				command: command,
				err:     err,
			}
		}

		if log != nil {
			defer log.Close()
		}

		if result == nil {
			result = RunCommand(command, log)
		}
		taskResult.commands[slot] = *result
		reporter.ReportResult(result)
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
