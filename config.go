package gofer

import (
	"io/ioutil"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type GoferConfig struct {
	Tasks map[string]GoferTask
}

func (c *GoferConfig) ToCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Gofer"
	app.Usage = "Your loyal task runner"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "explain, e",
			Usage: "Explain the execution plan, without taking any action",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Suppress command output when running tasks",
		},
	}

	commands := []cli.Command{}
	for name, task := range c.Tasks {
		commands = append(commands, task.ToCommand(name))
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
