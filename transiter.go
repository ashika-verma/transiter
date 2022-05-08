package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jamespfennell/transiter/internal/argsflag"
	"github.com/jamespfennell/transiter/internal/client"
	"github.com/jamespfennell/transiter/internal/server"
	"github.com/urfave/cli/v2"
)

func main() {
	argsMap := map[string]string{}
	app := &cli.App{
		Name:  "Transiter",
		Usage: "web service for transit data",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "addr",
				Aliases: []string{"a"},
				Usage:   "address of the Transiter server's gRPC admin API",
				Value:   "localhost:8083",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "delete",
				Usage: "delete a transit system",
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return fmt.Errorf("must provide the ID of the system to delete")
					}
					return clientAction(func(ctx context.Context, client *client.Client) error {
						return client.DeleteSystem(ctx, c.Args().Get(0))
					})(c)
				},
			},
			{
				Name:  "install",
				Usage: "install a transit system",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "interpret the second argument as a local file path",
						Value:   false,
					},
					&cli.BoolFlag{
						Name:  "template",
						Usage: "indicates that the input file is a Go template",
						Value: false,
					},
					&cli.BoolFlag{
						Name:    "update",
						Aliases: []string{"u"},
						Usage:   "if the system is already installed, update it with the provided config",
						Value:   false,
					},
					argsflag.NewCliFlag("arg", "", argsMap),
				},
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return fmt.Errorf("must provide the ID of the system to delete")
					}
					// TODO: pass the file name using --file and url using --url
					if c.Args().Len() == 1 {
						return fmt.Errorf("must provide a URL or file path for the transit system Yaml config")
					}
					args := client.InstallSystemArgs{
						SystemId:     c.Args().Get(0),
						ConfigPath:   c.Args().Get(1),
						IsFile:       c.Bool("file"),
						AllowUpdate:  c.Bool("update"),
						IsTemplate:   c.Bool("template") || c.IsSet("arg"),
						TemplateArgs: c.Value("arg").(map[string]string),
					}
					return clientAction(func(ctx context.Context, client *client.Client) error {
						return client.InstallSystem(ctx, args)
					})(c)
				},
			},
			{
				Name:  "list",
				Usage: "list all installed transit systems",
				Action: clientAction(func(ctx context.Context, client *client.Client) error {
					return client.ListSystems(ctx)
				}),
			},
			{
				Name:  "scheduler",
				Usage: "perform operations on the Transiter server scheduler",
				Subcommands: []*cli.Command{
					{
						Name:  "status",
						Usage: "get the list of periodic update tasks currently scheduled",
						Action: clientAction(func(ctx context.Context, client *client.Client) error {
							return client.SchedulerStatus(ctx)
						}),
					},
					{
						Name:  "refresh",
						Usage: "refresh the set of feed auto update tasks the scheduler is scheduling",
						Action: clientAction(func(ctx context.Context, client *client.Client) error {
							return client.RefreshScheduler(ctx)
						}),
					},
				},
			},
			{
				Name:  "server",
				Usage: "run a Transiter server",
				Action: func(c *cli.Context) error {
					return server.Run(c.String("postgres-addr"))
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "postgres-addr",
						Usage: "Postgres database address",
						Value: "localhost:5432",
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func clientAction(f func(ctx context.Context, client *client.Client) error) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		client, err := client.New(c.String("addr"))
		if err != nil {
			return err
		}
		defer client.Close()
		// TODO: parse the error to remove RPC references
		// For example when a yaml config url provided to install is incorrect
		return f(context.Background(), client)
	}
}
