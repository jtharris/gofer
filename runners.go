package main

import (
	"os/exec"

	"github.com/urfave/cli"
)

type Runner interface {
	Run() (string, error)
}

type explainRunner struct {
	command string
}

func (e explainRunner) Run() (string, error) {
	return e.command, nil
}

type shellRunner struct {
	command string
	quiet   bool
}

func (s shellRunner) Run() (string, error) {
	// TODO:  Support other shells here?
	out, err := exec.Command("bash", "-c", s.command).CombinedOutput()

	if s.quiet {
		out = nil
	}

	return string(out), err
}

func NewRunner(command string, context *cli.Context) Runner {
	if context.Parent().Bool("explain") {
		return explainRunner{
			command: command,
		}
	}

	return shellRunner{
		command: command,
		quiet:   context.Parent().Bool("quiet"),
	}
}
