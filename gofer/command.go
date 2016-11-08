package gofer

import (
	"io"
	"os/exec"
)

type GoferCommandResult struct {
	Command string
	Ran     bool
	Err     error
}

func RunCommand(command string, out io.Writer) *GoferCommandResult {
	// TODO:  Support other shells here?
	c := exec.Command("bash", "-c", command)
	c.Stdout = out
	c.Stderr = out

	return &GoferCommandResult{
		Command: command,
		Ran:     true,
		Err:     c.Run(),
	}
}
