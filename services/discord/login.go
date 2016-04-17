package discord

import (
	"bufio"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/codegangsta/cli"
	"github.com/kazokuco/disco/util"
	"os"
)

func (srv *Service) Login(c *cli.Context) {
	args := c.Args()
	if len(args) != 1 {
		log.Fatal("No bot specified!")
	}

	filename := args[0]
	brainFilename := util.BrainFilenameForBotFilename(filename)

	brain, err := util.LoadBrainFromFile(brainFilename)
	if err != nil {
		log.WithError(err).Fatal("Couldn't load brain")
	}

	store, err := brain.Get("service", "discord")
	if err != nil {
		log.WithError(err).Fatal("Couldn't get store")
	}

	st := store.(*Store)
	s := bufio.NewScanner(os.Stdin)

	// If a token isn't given in an environment variable, ask for one from STDIN
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		if st.Auth.Token != "" {
			fmt.Printf("Press ENTER to keep using: %s\n", st.Auth.Token)
		}

		fmt.Printf("Token: ")
		if !s.Scan() {
			os.Exit(0)
		}
		text := s.Text()
		if text != "" {
			token = text
			st.Auth.Token = text
		}
	}

	// Make a Discord session using it; no auth is checked here
	session, err := discordgo.New(token)
	if err != nil {
		log.WithError(err).Fatal("Couldn't connect to Discord")
	}

	// Verify that the token can actually connect
	resultCh := make(chan bool)
	session.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		resultCh <- true
	})
	session.AddHandler(func(s *discordgo.Session, event *discordgo.Disconnect) {
		resultCh <- false
	})
	if err = session.Open(); err != nil {
		log.WithError(err).Fatal("Couldn't open a connection")
	}

	// If we can't authorize, abort
	result := <-resultCh
	if !result {
		log.Fatal("Couldn't sign into Discord")
	}

	// Dump the bot's brain to a file
	if err = util.StoreBrainToFile(&brain, brainFilename); err != nil {
		log.WithError(err).Fatal("Couldn't store brain")
	}
}
