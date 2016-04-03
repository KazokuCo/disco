package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	_ "github.com/kazokuco/disco/jobs"
	_ "github.com/kazokuco/disco/services"
	"os"
)

// Tasks that run before the application is invoked.
func before(c *cli.Context) error {
	// Set logging verbosity
	log.SetLevel(log.InfoLevel)
	if c.GlobalBool("verbose") {
		log.SetLevel(log.DebugLevel)
	}

	return nil
}

func main() {
	cli.HelpFlag.Name = "help"
	cli.VersionFlag.Name = "version"

	app := cli.NewApp()
	app.Name = "disco"
	app.Usage = "A bot for the performance of menial tasks"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "Print debug information",
		},
	}
	app.Before = before
	app.Commands = []cli.Command{
		cli.Command{
			Name:      "run",
			Aliases:   []string{"r"},
			ArgsUsage: "bot.yml",
			Usage:     "Runs a bot",
			Action:    actionRun,
		},
		cli.Command{
			Name:      "login",
			ArgsUsage: "bot.yml service",
			Usage:     "Logs into the specified service",
			Action:    actionLogin,
		},
	}
	app.Run(os.Args)
}
