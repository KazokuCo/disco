package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/kazokuco/disco/bot"
)

func actionAvatar(c *cli.Context) {
	args := c.Args()

	if len(args) != 3 {
		log.Fatal("Usage: disco avatar bot service filename")
	}

	filename := args[0]
	brainFilename := BrainFilenameForBotFilename(filename)
	serviceName := args[1]
	avatarFilename := args[2]

	service := bot.GetService(serviceName)
	if service == nil {
		log.WithField("id", serviceName).Fatal("Unknown service")
	}

	brain, err := LoadBrainFromFile(brainFilename)
	if err != nil {
		log.WithError(err).Fatal("Couldn't load brain")
	}

	store, err := brain.Get(bot.TypeService, serviceName)
	if err != nil {
		log.WithError(err).Fatal("Couldn't get service store")
	}

	service.UpdateAvatar(store, avatarFilename)
}
