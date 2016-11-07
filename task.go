package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

type GoferTaskDefinition struct {
	Description string
	Parallel    bool
	Commands    []string
	Needs       []string
}

func (td *GoferTaskDefinition) ToGoferTask(name string) *GoferTask {
	return &GoferTask{
		Definition: td,
		Commands:   td.Commands, // TODO:  Expand macros here
		Name:       name,
	}
}

type GoferTask struct {
	Definition *GoferTaskDefinition
	Name       string
	TaskChain  []*GoferTask
	LogToFile  bool
	Commands   []string
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
		Usage: t.Definition.Description,
		Action: func(context *cli.Context) error {
			// Execute all dependencies before running the task
			// This includes the current task as well
			for _, task := range t.TaskChain {
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

	// TODO:  And once again - need a conversion step
	task.Commands = task.Definition.Commands

	if task.Definition.Parallel && task.LogToFile {
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
