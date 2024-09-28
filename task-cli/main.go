package main

import (
	"log"
	"os"

	"github.com/pbnjk/backend/task-cli/cmd"
)

func main() {
	app := cmd.New()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
