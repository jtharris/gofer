package main

import (
	"fmt"
	"os"

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
	LogToFile    bool
}

func (t GoferTask) CreateLogFile(slot int) (*os.File, error) {
	if !t.LogToFile {
		return os.Stdout, nil
	}

	err := os.MkdirAll("gofer-logs", os.ModePerm)

	if err != nil {
		return nil, err
	}

	fileName := fmt.Sprintf("gofer-logs/%s-%d.log", t.Name, slot)
	return os.Create(fileName)
}

func (t GoferTask) ToCommand() cli.Command {
	return cli.Command{
		Name:  t.Name,
		Usage: t.Description,
		Action: func(context *cli.Context) error {
			// Execute all dependencies before running the task
			// This includes the current task as well
			for _, task := range t.Dependencies {
				result := NewTaskRunner(context, task).Run()

				for _, r := range result.commands {
					if r.err != nil {
						return r.err
					}
				}
			}

			return nil
		},
	}
}

type GoferTaskResult struct {
	commands []GoferCommandResult
}

type GoferTaskRunner interface {
	Run() GoferTaskResult
}

func NewTaskRunner(context *cli.Context, task *GoferTask) GoferTaskRunner {
	if context.Parent().Bool("explain") {
		return explainTaskRunner{
			task: task,
		}
	}

	// TODO:  This is an odd place for this
	//        Make an explicit conversion step!
	task.LogToFile = !context.Parent().Bool("no-logs")

	if task.Parallel && task.LogToFile {
		return parallelTaskRunner{
			quiet: context.Parent().Bool("quiet"),
			task:  task,
		}
	}

	return serialTaskRunner{
		quiet: context.Parent().Bool("quiet"),
		task:  task,
	}
}
