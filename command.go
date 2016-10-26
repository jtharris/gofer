package main

import "os/exec"

type GoferCommandResult struct {
	command string
	output  string
	err     error
}

func RunCommand(command string) GoferCommandResult {
	// TODO:  Support other shells here?
	out, err := exec.Command("bash", "-c", command).CombinedOutput()

	return GoferCommandResult{
		command: command,
		output:  string(out),
		err:     err,
	}
}
