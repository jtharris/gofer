package main

import (
	"log"
	"os"
)

func main() {
	config, err := NewConfig("gofer.yml")

	if err != nil {
		log.Fatal(err)
	}

	app, err := config.ToCliApp()

	if err != nil {
		log.Fatal(err)
	}

	app.Run(os.Args)
}
