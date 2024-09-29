package cmd

import (
	"time"

	"github.com/urfave/cli/v2"
)

func New() *cli.App {
	return &cli.App{
		Name:     "expense-tracker",
		Usage:    "Expenses tracker",
		Version:  "1.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Pedro Buitrago",
				Email: "pedrobuitragons@gmail.com",
			},
		},
		Description: "Tool for helping you to track your day-to-day expenses",

		Before: Load,
		After:  Save,

		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Adds a new expense",
				Action:  HandleAdd,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "description",
						Aliases:     []string{"d"},
						Usage:       "The task's description",
						DefaultText: "an expense",
					},
					&cli.Float64Flag{
						Name:    "amount",
						Aliases: []string{"a"},
						Usage:   "The price of the expense",
					},
				},
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "Updates an expense",
				Action:  HandleUpdate,
				Flags: []cli.Flag{
					&cli.Uint64Flag{
						Name:  "id",
						Usage: "THe ID of the expense to update",
					},
					&cli.StringFlag{
						Name:    "description",
						Aliases: []string{"d"},
						Usage:   "The task's description",
					},
					&cli.Float64Flag{
						Name:    "amount",
						Aliases: []string{"a"},
						Usage:   "The price of the expense",
					},
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "Lists all expenses",
				Action:  HandleList,
			},
			{
				Name:    "summary",
				Aliases: []string{"s"},
				Usage:   "Prints a summary of all expenses",
				Action:  HandleSummary,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "month",
						Usage: "Limits summary to a given month",
					},
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "Deletes an expense",
				Action:  HandleDelete,
				Flags: []cli.Flag{
					&cli.Uint64Flag{
						Name:  "id",
						Usage: "The ID of the expense to delete",
					},
				},
			},
			{
				Name:    "set-limit",
				Aliases: []string{"sl"},
				Usage:   "Sets a monthly expense limit",
				Action:  HandleSetLimit,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "month",
						Usage: "Month to which the limit applies. If not set, limit applies to all months",
					},
					&cli.Float64Flag{
						Name:  "amount",
						Usage: "The monthly limit",
					},
				},
			},
			{
				Name:    "save",
				Aliases: []string{"sv"},
				Usage:   "Saves expenses to a CSV file",
				Action:  HandleSave,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "output",
						Aliases:     []string{"o"},
						Usage:       "Path to CSV file to output to",
						DefaultText: "output.csv",
					},
				},
			},
		},
	}
}
