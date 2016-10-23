package main

import (
	"log"

	"github.com/urfave/cli"
)

// TODO:  Support dependencies later
type GoferTask struct {
	Description string
	Commands    []string
	Needs       []string

	// TODO:  Is this mixing concerns?  YAML representation vs. internal
	// TODO:  Also consider renaming - perhaps "TaskChain" ?
	//        This list will include a reference to itself as well
	Dependencies []*GoferTask
}

func (t GoferTask) Execute(r Runner) error {
	for _, command := range t.Commands {
		output, err := r.Run(command)

		if len(output) > 0 {
			log.Println(output)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (t GoferTask) ToCommand(name string) cli.Command {
	return cli.Command{
		Name:  name,
		Usage: t.Description,
		Action: func(context *cli.Context) error {
			runner := NewRunner(context)
			// Execute all dependencies before running the task
			// This includes the current task as well
			for _, task := range t.Dependencies {
				err := task.Execute(runner)

				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}
