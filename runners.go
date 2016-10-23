package main

import (
	"os/exec"

	"github.com/urfave/cli"
)

type Runner interface {
	Run(command string) (string, error)
}

type explainRunner struct{}

func (e *explainRunner) Run(command string) (string, error) {
	return command, nil
}

type shellRunner struct {
	quiet bool
}

func (s *shellRunner) Run(command string) (string, error) {
	// TODO:  Support other shells here?
	out, err := exec.Command("bash", "-c", command).CombinedOutput()

	if s.quiet {
		out = nil
	}

	return string(out), err
}

func NewRunner(context *cli.Context) Runner {
	if context.Parent().Bool("explain") {
		return &explainRunner{}
	}

	return &shellRunner{
		quiet: context.Parent().Bool("quiet"),
	}
}
