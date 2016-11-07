package main

import (
	"log"
	"os"
)

func main() {
	configDefinition, err := NewConfigDefinition("gofer.yml")

	if err != nil {
		log.Fatal(err)
	}

	config, err := configDefinition.ToConfig()

	if err != nil {
		log.Fatal(err)
	}

	config.ToCliApp().Run(os.Args)
}
