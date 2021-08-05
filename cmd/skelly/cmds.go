package main

import (
	"github.com/davidvader/skelly/router"
	"github.com/davidvader/skelly/skelly"
	"github.com/davidvader/skelly/util"
	"github.com/urfave/cli/v2"
)

// commands is a collection of actions available via the CLI
var (
	// serverCmd defines the command for running the http server.
	serverCmd = &cli.Command{
		Name:        "server",
		Category:    "Server",
		Aliases:     []string{"s"},
		Description: "Use this command to run the http server.",
		Usage:       "Run the Vela Slack bot http server",
		Action:      server,
		Flags: []cli.Flag{
			&cli.StringFlag{
				EnvVars: []string{"SKELLY_PORT"},
				Name:    "port",
				Usage:   "port for Skelly server",
				Value:   "8080",
			},
		},
	}

	// reactionCmds defines the main command for controlling reactions.
	// trigger defines the command for simulating a skelly a slack reaction.
	reactionCmds = []*cli.Command{
		{
			Name:        "reaction",
			Category:    "Reaction",
			Aliases:     []string{"r"},
			Description: "Use this command to control reactions for a specified channel.",
			Usage:       "Controls reactions for a specified channel",
			Subcommands: []*cli.Command{
				{
					Name:        "view",
					Category:    "Reaction",
					Aliases:     []string{"v", "get"},
					Description: "Use this command to view a reaction for a specified channel.",
					Usage:       "View reaction for a specified channel",
					Before:      validateView,
					Action:      view,
					Flags: []cli.Flag{
						&cli.StringFlag{
							EnvVars: []string{"SKELLY_CHANNEL"},
							Name:    "channel",
							Aliases: []string{"c"},
							Usage:   "for which channel to retrieve a reaction",
							Value:   "",
						},
					},
				},
				{
					Name:        "list",
					Category:    "Reaction",
					Aliases:     []string{"l"},
					Description: "Use this command to list reactions for a specified channel.",
					Usage:       "List reactions for a specified channel",
					Before:      validateList,
					Action:      list,
					Flags: []cli.Flag{
						&cli.StringFlag{
							EnvVars: []string{"SKELLY_CHANNEL"},
							Name:    "channel",
							Aliases: []string{"c"},
							Usage:   "for which channel to retrieve reactions",
							Value:   "",
						},
					},
				},
				{
					Name:        "clear",
					Category:    "Reaction",
					Aliases:     []string{"c"},
					Description: "Use this command to clear reactions for a specified channel.",
					Usage:       "Clear reactions for a specified channel",
					Before:      validateClear,
					Action:      clear,
					Flags: []cli.Flag{
						&cli.StringFlag{
							EnvVars: []string{"SKELLY_CHANNEL"},
							Name:    "channel",
							Aliases: []string{"c"},
							Usage:   "for which channel to clear reactions",
							Value:   "",
						},
					},
				},
				{
					Name:        "add",
					Category:    "Reaction",
					Aliases:     []string{"a"},
					Description: "Use this command to add a reaction for a specified channel.",
					Usage:       "Add a reaction for a specified channel",
					Before:      validateAdd,
					Action:      add,
					Flags: []cli.Flag{
						&cli.StringFlag{
							EnvVars: []string{"SKELLY_CHANNEL"},
							Name:    "channel",
							Aliases: []string{"c"},
							Usage:   "for which channel to add",
							Value:   "",
						},
						&cli.StringFlag{
							Name:    "response",
							Aliases: []string{"r"},
							Usage:   "what message to respond with",
							Value:   "",
						},
					},
				},
				{
					Name:        "update",
					Category:    "Reaction",
					Aliases:     []string{"u"},
					Description: "Use this command to update a reaction for a specified channel.",
					Usage:       "Update a reaction for a specified channel",
					Before:      validateUpdate,
					Action:      update,
					Flags: []cli.Flag{
						&cli.StringFlag{
							EnvVars: []string{"SKELLY_CHANNEL"},
							Name:    "channel",
							Aliases: []string{"c"},
							Usage:   "for which channel to update",
							Value:   "",
						},
						&cli.StringFlag{
							Name:    "response",
							Aliases: []string{"r"},
							Usage:   "what message to respond with",
							Value:   "",
						},
					},
				},
				{
					Name:        "delete",
					Category:    "Reaction",
					Aliases:     []string{"d"},
					Description: "Use this command to delete reactions for a specified channel.",
					Usage:       "Delete reactions for a specified channel",
					Before:      validateDelete,
					Action:      delete,
					Flags: []cli.Flag{
						&cli.StringFlag{
							EnvVars: []string{"SKELLY_CHANNEL"},
							Name:    "channel",
							Aliases: []string{"c"},
							Usage:   "for which channel to delete",
							Value:   "",
						},
					},
				},
				{
					Name:        "trigger",
					Category:    "Reaction",
					Aliases:     []string{"t"},
					Description: "Use this command to simulate a reaction for a specified channel (and ts).",
					Usage:       "Trigger a reaction for a specified channel (and ts)",
					Before:      validateTrigger,
					Action:      trigger,
					Flags: []cli.Flag{
						&cli.StringFlag{
							EnvVars: []string{"SKELLY_CHANNEL"},
							Name:    "channel",
							Aliases: []string{"c"},
							Usage:   "which channel to trigger a reaction in",
							Value:   "",
						},
						&cli.StringFlag{
							EnvVars: []string{"SKELLY_USER"},
							Name:    "user",
							Aliases: []string{"u"},
							Usage:   "which user to trigger a reaction on",
							Value:   "",
						},
						&cli.StringFlag{
							EnvVars: []string{"SKELLY_TIMESTAMP"},
							Name:    "ts",
							Usage:   "which message timestamp to trigger a reaction on",
							Value:   "none",
						},
					},
				},
			},
		},
	}
)

func cmds() []*cli.Command {
	return append(reactionCmds, serverCmd)
}

// validateView is a helper function to load global configuration if set
// via config or environment and validate the user input in the command
func validateView(c *cli.Context) error {
	// validate the user input in the command
	if len(c.String("channel")) == 0 {
		return util.InvalidCommand("channel")
	}

	return nil
}

// validateList is a helper function to load global configuration if set
// via config or environment and validate the user input in the command
func validateList(c *cli.Context) error {

	// validate the user input in the command
	if len(c.String("channel")) == 0 {
		return util.InvalidCommand("channel")
	}

	return nil
}

// validateClear is a helper function to load global configuration if set
// via config or environment and validate the user input in the command
func validateClear(c *cli.Context) error {

	// validate the user input in the command
	if len(c.String("channel")) == 0 {
		return util.InvalidCommand("channel")
	}

	return nil
}

// validateAdd is a helper function to load global configuration if set
// via config or environment and validate the user input in the command
func validateAdd(c *cli.Context) error {

	// validate the user input in the command
	if len(c.String("channel")) == 0 {
		return util.InvalidCommand("channel")
	}
	if len(c.String("response")) == 0 {
		return util.InvalidCommand("response")
	}

	return nil
}

// validateUpdate is a helper function to load global configuration if set
// via config or environment and validate the user input in the command
func validateUpdate(c *cli.Context) error {

	// validate the user input in the command
	if len(c.String("channel")) == 0 {
		return util.InvalidCommand("channel")
	}
	if len(c.String("response")) == 0 {
		return util.InvalidCommand("response")
	}

	return nil
}

// validateDelete is a helper function to load global configuration if set
// via config or environment and validate the user input in the command
func validateDelete(c *cli.Context) error {

	// validate the user input in the command
	if len(c.String("channel")) == 0 {
		return util.InvalidCommand("channel")
	}

	return nil
}

// validateTrigger is a helper function to load global configuration if set
// via config or environment and validate the user input in the command
func validateTrigger(c *cli.Context) error {

	// validate the user input in the command
	if len(c.String("channel")) == 0 {
		return util.InvalidCommand("channel")
	}

	return nil
}

// validateSkellyStats is a helper function to load global configuration if set
// via config or environment and validate the user input in the command
func validateSkellyStats(c *cli.Context) error {

	return nil
}

// server is a wrapper around running router.Run via the CLI
func server(c *cli.Context) error {
	return router.Run(c.String("port"))
}

// view is a wrapper around running skelly.View via the CLI
func view(c *cli.Context) error {
	return skelly.View(c.String("channel"))
}

// list is a wrapper around running skelly.List via the CLI
func list(c *cli.Context) error {
	return skelly.List(c.String("token"), c.String("channel"))
}

// clear is a wrapper around running skelly.List via the CLI
func clear(c *cli.Context) error {
	return skelly.Clear(c.String("channel"))
}

// add is a wrapper around running skelly.Add via the CLI
func add(c *cli.Context) error {
	return skelly.Add(c.String("token"), c.String("channel"), c.String("response"))
}

// update is a wrapper around running skelly.Update via the CLI
func update(c *cli.Context) error {
	return skelly.Update(c.String("token"), c.String("channel"), c.String("response"))
}

// delete is a wrapper around running skelly.Delete via the CLI
func delete(c *cli.Context) error {
	return skelly.Delete(c.String("token"), c.String("channel"))
}

// trigger is a wrapper around running skelly.Trigger via the CLI
func trigger(c *cli.Context) error {
	return skelly.Trigger(c.String("token"), c.String("channel"), c.String("user"), c.String("ts"))
}
