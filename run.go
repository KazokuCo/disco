package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/kazokuco/disco/bot"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func actionRun(c *cli.Context) {
	args := c.Args()

	if len(args) != 1 {
		log.Fatal("No bot specified!")
	}

	filename := args[0]
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithError(err).Fatal("Couldn't open file")
	}

	b := bot.New()
	err = yaml.Unmarshal(data, &b)
	if err != nil {
		log.WithError(err).Fatal("Couldn't load config")
	}

	b.Run()
}
