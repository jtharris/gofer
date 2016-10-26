package main

import (
	"os/exec"

	"github.com/urfave/cli"
)

type CommandRunner struct {
	quiet bool
}

func (c CommandRunner) Run(command string) (string, error) {
	// TODO:  Support other shells here?
	out, err := exec.Command("bash", "-c", command).CombinedOutput()

	if c.quiet {
		out = nil
	}

	return string(out), err
}

// TODO:  This abstraction might not be needed  now
//        Consider simplifying this
func NewCommandRunner(context *cli.Context) CommandRunner {
	return CommandRunner{
		quiet: context.Parent().Bool("quiet"),
	}
}
