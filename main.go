package main

import (
	"gofer/gofer"
	"log"
	"os"
)

func main() {
	configDefinition, err := gofer.NewConfigDefinition("gofer.yml")

	if err != nil {
		log.Fatal(err)
	}

	config, err := configDefinition.ToConfig()

	if err != nil {
		log.Fatal(err)
	}

	config.ToCliApp().Run(os.Args)
}
