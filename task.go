package gofer

import (
	"log"
	"os/exec"

	"github.com/urfave/cli"
)

// TODO:  Support dependencies later
type GoferTask struct {
	Description string
	Command     string
}

func (t GoferTask) ToCommand(name string) cli.Command {
	return cli.Command{
		Name:  name,
		Usage: t.Description,
		Action: func(c *cli.Context) error {
			if c.Parent().Bool("explain") {
				log.Println("Command:", t.Command)
				return nil
			}

			out, err := exec.Command("bash", "-c", t.Command).Output()

			if c.Parent().Bool("quiet") == false {
				log.Println(string(out))
			}

			return err
		},
	}
}
