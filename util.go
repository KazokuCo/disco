package main

import (
	"fmt"
	"github.com/kazokuco/disco/bot"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func BrainFilenameForBotFilename(filename string) string {
	return fmt.Sprintf("%s.brain", filename)
}

func LoadBotFromFile(filename string) (b bot.Bot, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return b, err
	}

	b = bot.New()
	if err = yaml.Unmarshal(data, &b); err != nil {
		return b, err
	}

	return b, nil
}

func LoadBrainFromFile(filename string) (brain bot.Brain, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return brain, err
	}

	brain = bot.NewBrain()
	if err = yaml.Unmarshal(data, &brain); err != nil {
		return brain, err
	}

	return brain, nil
}

func StoreBrainToFile(brain *bot.Brain, filename string) error {
	data, err := yaml.Marshal(brain)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}
