package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/kazokuco/disco/bot"
)

func actionLogin(c *cli.Context) {
	args := c.Args()

	if len(args) != 2 {
		log.Fatal("Usage: disco login bot service")
	}

	filename := args[0]
	brainFilename := BrainFilenameForBotFilename(filename)
	serviceName := args[1]

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

	if persist := service.Login(store); persist {
		if err = StoreBrainToFile(&brain, brainFilename); err != nil {
			log.WithError(err).Fatal("Couldn't store brain")
		}
	}
}
