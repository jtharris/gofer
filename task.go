package main

import (
	"log"

	"github.com/urfave/cli"
)

type GoferTask struct {
	Description string
	Parallel    bool
	Commands    []string
	Needs       []string

	// TODO:  Is this mixing concerns?  YAML representation vs. internal
	// TODO:  Also consider renaming - perhaps "TaskChain" ?
	//        This list will include a reference to itself as well
	Name         string // this should be in a separate struct - make YAML representation different
	Dependencies []*GoferTask
}

type GoferTaskResult struct {
	Output string
	Error  error
}

type GoferTaskRunner interface {
	Run() []GoferTaskResult
}

type serialTaskRunner struct {
	commandRunner CommandRunner
	task          *GoferTask
}

func (r serialTaskRunner) Run() []GoferTaskResult {
	results := make([]GoferTaskResult, len(r.task.Commands))
	for i, command := range r.task.Commands {
		output, err := r.commandRunner.Run(command)

		results[i] = GoferTaskResult{
			Output: output,
			Error:  err,
		}

		if len(output) > 0 {
			log.Println(output)
		}

		if err != nil {
			return results
		}
	}

	return results
}

type explainTaskRunner struct {
	task *GoferTask
}

func (r explainTaskRunner) Run() []GoferTaskResult {
	log.Println(r.task.Name)
	for _, command := range r.task.Commands {
		log.Println("  |>", command)
	}

	return []GoferTaskResult{}
}

func NewTaskRunner(context *cli.Context, task *GoferTask) GoferTaskRunner {
	if context.Parent().Bool("explain") {
		return explainTaskRunner{
			task: task,
		}
	}

	return serialTaskRunner{
		commandRunner: NewCommandRunner(context),
		task:          task,
	}
}

func (t GoferTask) ToCommand() cli.Command {
	return cli.Command{
		Name:  t.Name,
		Usage: t.Description,
		Action: func(context *cli.Context) error {
			// Execute all dependencies before running the task
			// This includes the current task as well
			for _, task := range t.Dependencies {
				results := NewTaskRunner(context, task).Run()

				for _, result := range results {
					if result.Error != nil {
						return result.Error
					}
				}
			}

			return nil
		},
	}
}
