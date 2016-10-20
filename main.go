package gofer

import (
	"log"
	"os"
)

func main() {
	config, err := NewConfig("gofer.yml")

	if err != nil {
		log.Fatal(err)
	}

	config.ToCliApp().Run(os.Args)
}
