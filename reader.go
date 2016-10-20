package main

import (
	"io/ioutil"
	"log"
	"os/exec"

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

// TODO:  Support dependencies later
type GoferTask struct {
	Description string
	Command     string
}

func (t GoferTask) ToCommand(name string) cli.Command {
	return cli.Command{
		Name:  name,
		Usage: t.Description,
		Action: func(c *cli.Context) error {
			if c.Parent().Bool("explain") {
				log.Println("Command:", t.Command)
				return nil
			}

			out, err := exec.Command("bash", "-c", t.Command).Output()

			// TODO:  Introduce verbose mode
			// TODO:  Introduce quiet mode
			// TODO:  Introduce file logging
			log.Println(string(out))

			return err
		},
	}
}