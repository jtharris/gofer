package main

import (
	"fmt"

	"github.com/fatih/color"
)

type CommandReporter struct {
	quiet    bool
	parallel bool
}

func (r *CommandReporter) ReportTaskStart(task *GoferTask) {
	color.New(color.FgHiYellow, color.Bold).Println(task.Name)
}

func (r *CommandReporter) ReportResult(result *GoferCommandResult) {
	if r.quiet {
		return
	}
	var preCommandString string

	if r.parallel {
		preCommand := color.New(color.FgGreen, color.Bold).SprintFunc()
		preCommandString = preCommand("  |")
	} else {
		preCommand := color.New(color.FgWhite, color.Bold).SprintFunc()
		preCommandString = preCommand("  |>")
	}

	if len(result.output) == 0 {
		fmt.Println(preCommandString, result.command)
	} else if result.err != nil {
		fmt.Println(preCommandString, result.command, color.RedString("\u2717"))
		fmt.Println("   ", result.output)
	} else {
		fmt.Println(preCommandString, result.command, color.GreenString("\u2713"))
	}
}
