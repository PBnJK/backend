package cmd

import (
	"time"

	"github.com/urfave/cli/v2"
)

func New() *cli.App {
	return &cli.App{
		Name:     "task-cli",
		Usage:    "CLI task manager",
		Version:  "1.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Pedro Buitrago",
				Email: "pedrobuitragons@gmail.com",
			},
		},
		Description: "Tool for helping you to manage tasks using todo lists",

		Before: Load,
		After:  Save,

		Commands: []*cli.Command{
			{
				Name:      "add",
				Aliases:   []string{"a"},
				Usage:     "Adds a new task",
				UsageText: "task-cli [add, a] [task name]",
				Action:    HandleAdd,
			},
			{
				Name:      "update",
				Aliases:   []string{"u"},
				Usage:     "Updates a task",
				UsageText: "task-cli [update, u] [task id] [task name]",
				Action:    HandleUpdate,
			},
			{
				Name:      "delete",
				Aliases:   []string{"d"},
				Usage:     "Deletes a task",
				UsageText: "task-cli [delete, d] <task id or status>",
				Action:    HandleDelete,
				Subcommands: []*cli.Command{
					{
						Name:     "done",
						Aliases:  []string{"d"},
						Usage:    "Deletes all completed tasks",
						Action:   HandleDeleteDone,
						Category: "list",
					},
					{
						Name:     "todo",
						Aliases:  []string{"t"},
						Usage:    "Deletes all tasks that are yet to be started",
						Action:   HandleDeleteTodo,
						Category: "list",
					},
					{
						Name:     "in-progress",
						Aliases:  []string{"p"},
						Usage:    "Deletes all in-progress tasks",
						Action:   HandleDeleteInProgress,
						Category: "list",
					},
				},
			},
			{
				Name:      "mark-in-progress",
				Aliases:   []string{"mp"},
				Usage:     "Marks a task as in-progress",
				UsageText: "task-cli [mark-in-progress, mp] [task id]",
				Action:    HandleMarkInProgress,
			},
			{
				Name:      "mark-done",
				Aliases:   []string{"md"},
				Usage:     "Marks a task as done",
				UsageText: "task-cli [mark-done, md] [task id]",
				Action:    HandleMarkDone,
			},
			{
				Name:      "list",
				Aliases:   []string{"l"},
				Usage:     "Lists all tasks",
				UsageText: "task-cli [list, l] <type>",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "verbose",
						Aliases: []string{"v"},
					},
				},
				Action: HandleList,
				Subcommands: []*cli.Command{
					{
						Name:     "done",
						Aliases:  []string{"d"},
						Usage:    "Lists all completed tasks",
						Action:   HandleListDone,
						Category: "list", Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "verbose",
								Aliases: []string{"v"},
							},
						},
					},
					{
						Name:     "todo",
						Aliases:  []string{"t"},
						Usage:    "Lists all tasks that are yet to be started",
						Action:   HandleListTodo,
						Category: "list", Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "verbose",
								Aliases: []string{"v"},
							},
						},
					},
					{
						Name:     "in-progress",
						Aliases:  []string{"p"},
						Usage:    "Lists all in-progress tasks",
						Action:   HandleListInProgress,
						Category: "list", Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "verbose",
								Aliases: []string{"v"},
							},
						},
					},
				},
			},
		},
	}
}
