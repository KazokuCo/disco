package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/kazokuco/disco/bot"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func actionRun(c *cli.Context) {
	args := c.Args()

	if len(args) != 1 {
		log.Fatal("No bot specified!")
	}

	filename := args[0]
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithError(err).Fatal("Couldn't load bot")
	}

	b := bot.New()
	if err = yaml.Unmarshal(data, &b); err != nil {
		log.WithError(err).Fatal("Couldn't load bot")
	}

	brainFilename := fmt.Sprintf("%s.brain", filename)
	brainData, err := ioutil.ReadFile(brainFilename)
	if err != nil && !os.IsNotExist(err) {
		log.WithError(err).Fatal("Couldn't load brain")
	}

	brain := bot.NewBrain()
	if err = yaml.Unmarshal(brainData, &brain); err != nil {
		log.WithError(err).Fatal("Couldn't load brain")
	}

	err = b.Run(&brain)
	if err != nil {
		log.WithError(err).Fatal("Error")
	}

	brainData, err = yaml.Marshal(brain)
	ioutil.WriteFile(brainFilename, brainData, 0644)
}
