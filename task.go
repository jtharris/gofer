package main

import (
	"log"
	"os/exec"

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
		Action: func(c *cli.Context) error {
			explain := c.Parent().Bool("explain")
			quiet := !explain && c.Parent().Bool("quiet")

			for i, command := range t.Commands {
				if explain {
					log.Println("Command", i, ": ", command)
				} else {
					out, err := exec.Command("bash", "-c", command).Output()

					if !quiet {
						log.Println(string(out))
					}

					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}
