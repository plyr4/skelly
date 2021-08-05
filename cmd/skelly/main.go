package main

import (
	"fmt"
	"os"
	"time"

	"github.com/davidvader/skelly/db"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {

	// load environment from .env
	godotenv.Load(".env")

	app := cli.NewApp()

	// App Information
	app.Name = "skelly"
	app.HelpName = "skelly"
	app.Usage = "Slack bot for automatically reacting to typing, written in Go"
	app.Authors = []*cli.Author{
		{
			Name: "David Vader",
		},
	}

	// App Metadata
	app.Compiled = time.Now()
	app.EnableBashCompletion = true
	app.Version = "0.1.0"

	// App Commands
	app.Commands = cmds()

	// App Flags
	app.Flags = []cli.Flag{

		&cli.StringFlag{
			EnvVars: []string{"LOG_LEVEL"},
			Name:    "log-level",
			Aliases: []string{"log"},
			Usage:   "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:   "info",
		},
		&cli.StringFlag{
			EnvVars: []string{"SKELLY_BOT_TOKEN"},
			Name:    "token",
			Aliases: []string{"t"},
			Usage:   "bot token for interacting with the Slack API. See: https://api.slack.com/authentication/basics",
			Value:   "",
		},
		&cli.BoolFlag{
			EnvVars: []string{"SKIP_LOAD"},
			Name:    "skip-load",
			Aliases: []string{"sl"},
			Usage:   "set whether to skip load - options: (true|false)",
			Value:   false,
		},
	}

	// App Configurations
	app.Before = load

	// verify the database config
	err := db.Verify()
	if err != nil {
		panic(err)
	}

	// Run App
	err = app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

// load is a helper function that loads the necessary configuration for the CLI.
func load(c *cli.Context) error {

	// set log level based on config
	switch c.String("log-level") {
	case "t", "trace", "Trace", "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "i", "info", "Info", "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "w", "warn", "Warn", "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error", "Error", "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal", "Fatal", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic", "Panic", "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	}

	// validate configurations
	err := validate(c)
	if err != nil {
		return err
	}

	return nil
}

// validate is a helper function that ensures the necessary configuration is set for the CLI.
func validate(c *cli.Context) error {
	args := c.Args()

	// skip validate if help argument is provided
	for _, arg := range args.Slice() {
		if arg == "--help" || arg == "-h" {
			return nil
		}
	}

	// additional validations would go here
	if len(c.String("token")) == 0 {
		return fmt.Errorf("no bot token provided")
	}

	return nil
}
