package main

import (
	"io"
	"os/exec"
)

type GoferCommandResult struct {
	command string
	ran     bool
	err     error
}

func RunCommand(command string, out io.Writer) *GoferCommandResult {
	// TODO:  Support other shells here?
	c := exec.Command("bash", "-c", command)
	c.Stdout = out
	c.Stderr = out

	return &GoferCommandResult{
		command: command,
		ran:     true,
		err:     c.Run(),
	}
}
