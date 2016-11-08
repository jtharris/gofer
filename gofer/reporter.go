package gofer

import (
	"fmt"

	"github.com/fatih/color"
)

type CommandReporter struct {
	Quiet    bool
	Parallel bool
}

func (r *CommandReporter) ReportTaskStart(task *GoferTask) {
	color.New(color.FgHiYellow, color.Bold).Println(task.Name)
}

func (r *CommandReporter) ReportResult(result *GoferCommandResult) {
	if r.Quiet {
		return
	}
	var preCommandString string

	if r.Parallel {
		preCommand := color.New(color.FgGreen, color.Bold).SprintFunc()
		preCommandString = preCommand("  |")
	} else {
		preCommand := color.New(color.FgWhite, color.Bold).SprintFunc()
		preCommandString = preCommand("  |>")
	}

	if !result.Ran {
		fmt.Println(preCommandString, result.Command)
	} else if result.Err != nil {
		fmt.Println(preCommandString, result.Command, color.RedString("\u2717"))
	} else {
		fmt.Println(preCommandString, result.Command, color.GreenString("\u2713"))
	}
}
