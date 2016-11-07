package main

import (
	"io/ioutil"

	"github.com/stevenle/topsort"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type GoferConfigDefinition struct {
	Tasks map[string]*GoferTaskDefinition
}

func (gcd *GoferConfigDefinition) ToConfig() (*GoferConfig, error) {
	tasks := make(map[string]*GoferTask)
	for name, taskDef := range gcd.Tasks {
		tasks[name] = taskDef.ToGoferTask(name)
	}

	config := &GoferConfig{
		Tasks: tasks,
	}

	return config, config.populateDependentTasks()
}

type GoferConfig struct {
	Tasks map[string]*GoferTask
}

func (c *GoferConfig) ToCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Gofer"
	app.Usage = "Your loyal task runner"
	app.Version = "0.0.1"
	app.Flags = c.getGlobalFlags()
	app.Commands = c.getCLICommands()

	return app
}

func (c *GoferConfig) getGlobalFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "explain, e",
			Usage: "Explain the execution plan, without taking any action",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Suppress command output when running tasks",
		},
		cli.BoolFlag{
			Name:  "no-logs, l",
			Usage: "Stream task output to stdout, instead of writing to logs.  All tasks will run serially.",
		},
	}
}

func (c *GoferConfig) getCLICommands() []cli.Command {
	commands := make([]cli.Command, 0, len(c.Tasks))
	for _, task := range c.Tasks {
		commands = append(commands, task.ToCommand())
	}

	return commands
}

func (c *GoferConfig) populateDependentTasks() error {
	graph := topsort.NewGraph()

	// Build up the graph nodes
	for name := range c.Tasks {
		graph.AddNode(name)
	}

	// Add graph edges representing dependencies
	for name, task := range c.Tasks {
		for _, dependency := range task.Definition.Needs {
			graph.AddEdge(name, dependency)
		}
	}

	// Set the transitive dependencies
	for name, task := range c.Tasks {
		depNames, err := graph.TopSort(name)

		if err != nil {
			return err
		}

		task.TaskChain = []*GoferTask{}
		for _, depName := range depNames {
			task.TaskChain = append(task.TaskChain, c.Tasks[depName])
		}
	}

	return nil
}

func NewConfigDefinition(configFile string) (*GoferConfigDefinition, error) {
	definition := &GoferConfigDefinition{}
	configData, err := ioutil.ReadFile(configFile)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(configData, definition)

	if err != nil {
		return nil, err
	}

	return definition, nil
}
