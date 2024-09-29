package cmd

import (
	"time"

	"github.com/urfave/cli/v2"
)

func New() *cli.App {
	return &cli.App{
		Name:     "github-activity",
		Usage:    "Github activity fetcher",
		Version:  "1.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Pedro Buitrago",
				Email: "pedrobuitragons@gmail.com",
			},
		},
		Description: "Tool for printing out a user's Github activity",
		Action:      Handle,
	}
}
