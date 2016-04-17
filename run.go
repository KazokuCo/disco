package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/kazokuco/disco/util"
	"os"
	"os/signal"
)

func actionRun(c *cli.Context) {
	args := c.Args()

	if len(args) != 1 {
		log.Fatal("No bot specified!")
	}
	filename := args[0]
	brainFilename := util.BrainFilenameForBotFilename(filename)

	// Load the bot configuration
	b, err := util.LoadBotFromFile(filename)
	if err != nil {
		log.WithError(err).Fatal("Couldn't load bot")
	}

	// Load the bot's brain
	brain, err := util.LoadBrainFromFile(brainFilename)
	if err != nil {
		log.WithError(err).Fatal("Couldn't load brain")
	}

	// Listen for SIGINT (Ctrl+C)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	// Run the bot
	stop := make(chan interface{})
	go func() {
		if err = b.Run(&brain, stop); err != nil {
			log.WithError(err).Fatal("Error")
		}
	}()

	// Keep running until we get a SIGINT
	<-ch
	close(stop)

	log.Info("Shutting down")

	// Dump the bot's brain to a file
	if err = util.StoreBrainToFile(&brain, brainFilename); err != nil {
		log.WithError(err).Fatal("Couldn't store brain")
	}
}
