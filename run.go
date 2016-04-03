package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func actionRun(c *cli.Context) {
	args := c.Args()

	if len(args) != 1 {
		log.Fatal("No bot specified!")
	}
	filename := args[0]
	brainFilename := BrainFilenameForBotFilename(filename)

	b, err := LoadBotFromFile(filename)
	if err != nil {
		log.WithError(err).Fatal("Couldn't load bot")
	}

	brain, err := LoadBrainFromFile(brainFilename)
	if err != nil {
		log.WithError(err).Fatal("Couldn't load brain")
	}

	if err = b.Run(&brain); err != nil {
		log.WithError(err).Fatal("Error")
	}

	if err = StoreBrainToFile(&brain, brainFilename); err != nil {
		log.WithError(err).Fatal("Couldn't store brain")
	}
}
