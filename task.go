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
			return t.Execute(runner)
		},
	}
}
