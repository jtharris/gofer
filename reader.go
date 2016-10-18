package main

import (
	"fmt"
	"io/ioutil"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type GoferConfig struct {
	Name        string
	Description string
	Tasks       []GoferTask
}

func (c *GoferConfig) ToCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = c.Name
	app.Usage = c.Description

	commands := make([]cli.Command, len(c.Tasks))
	for i, task := range c.Tasks {
		commands[i] = task.ToCommand()
	}

	app.Commands = commands

	return app
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

// TODO:  Support dependencies later
type GoferTask struct {
	Name        string
	Description string
	Command     string
}

func (t GoferTask) ToCommand() cli.Command {
	return cli.Command{
		Name:  t.Name,
		Usage: t.Description,
		Action: func(c *cli.Context) error {
			fmt.Println("Running command:", t.Command)
			return nil
		},
	}
}
