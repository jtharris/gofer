package main

import (
	"log"

	"github.com/urfave/cli"
)

// TODO:  Support dependencies later
type GoferTask struct {
	Description string
	Commands    []string
}

func (t GoferTask) ToCommand(name string) cli.Command {
	return cli.Command{
		Name:  name,
		Usage: t.Description,
		Action: func(context *cli.Context) error {
			for _, command := range t.Commands {
				output, err := NewRunner(command, context).Run()

				if len(output) > 0 {
					log.Println(output)
				}

				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}
