package gofer

import "sync"

type serialTaskRunner struct {
	quiet bool
	task  *GoferTask
}

func (r serialTaskRunner) Run() GoferTaskResult {
	reporter := CommandReporter{
		Quiet:    r.quiet,
		Parallel: false,
	}
	taskResult := GoferTaskResult{
		Commands: make([]GoferCommandResult, len(r.task.Commands)),
	}

	reporter.ReportTaskStart(r.task)

	for i, command := range r.task.Commands {
		log, err := r.task.CreateLogFile(i)
		if err != nil {
			taskResult.Commands[i] = GoferCommandResult{
				Command: command,
				Err:     err,
			}

			break
		}

		if log != nil {
			mode, err := log.Stat()

			if err != nil && mode.Mode().IsRegular() {
				defer log.Close()
			}
		}

		result := RunCommand(command, log)
		taskResult.Commands[i] = *result
		reporter.ReportResult(result)

		if result.Err != nil {
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
		Quiet:    r.quiet,
		Parallel: true,
	}
	taskResult := GoferTaskResult{
		Commands: make([]GoferCommandResult, len(r.task.Commands)),
	}

	reporter.ReportTaskStart(r.task)

	run := func(slot int, command string) {
		defer wg.Done()
		var result *GoferCommandResult

		// TODO:  DRY this up w/ the serialTaskRunner - lots of copy-paste
		log, err := r.task.CreateLogFile(slot)
		if err != nil {
			result = &GoferCommandResult{
				Command: command,
				Err:     err,
			}
		}

		if log != nil {
			defer log.Close()
		}

		if result == nil {
			result = RunCommand(command, log)
		}
		taskResult.Commands[slot] = *result
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
		Quiet:    false,
		Parallel: r.task.Definition.Parallel,
	}

	reporter.ReportTaskStart(r.task)

	for _, command := range r.task.Commands {
		result := &GoferCommandResult{
			Command: command,
		}
		reporter.ReportResult(result)
	}

	return GoferTaskResult{}
}
