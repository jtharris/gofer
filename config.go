package main

import (
	"io/ioutil"

	"github.com/stevenle/topsort"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type GoferConfig struct {
	Tasks map[string]*GoferTask
}

func (c *GoferConfig) ToCliApp() (*cli.App, error) {
	app := cli.NewApp()
	app.Name = "Gofer"
	app.Usage = "Your loyal task runner"
	app.Version = "0.0.1"

	err := c.populateDependentTasks()

	app.Flags = c.getGlobalFlags()
	app.Commands = c.getCLICommands()

	return app, err
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
	}
}

func (c *GoferConfig) getCLICommands() []cli.Command {
	commands := []cli.Command{}
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
		// TODO:  Make this a different function altogether...
		task.Name = name

		for _, dependency := range task.Needs {
			graph.AddEdge(name, dependency)
		}
	}

	// Set the transitive dependencies
	for name, task := range c.Tasks {
		depNames, err := graph.TopSort(name)

		if err != nil {
			return err
		}

		task.Dependencies = []*GoferTask{}
		for _, depName := range depNames {
			task.Dependencies = append(task.Dependencies, c.Tasks[depName])
		}
	}

	return nil
}

func NewConfig(configFile string) (*GoferConfig, error) {
	config := &GoferConfig{}
	configData, err := ioutil.ReadFile(configFile)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(configData, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
