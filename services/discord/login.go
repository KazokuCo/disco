package discord

import (
	"bufio"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/kazokuco/disco/bot"
	"os"
)

func (srv *Service) Login(store bot.Store) bool {
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
			return false
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
		log.WithError(err).Error("Couldn't connect to Discord")
		return false
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
		log.WithError(err).Error("Couldn't open a connection")
		return false
	}

	// If we can't authorize, abort
	result := <-resultCh
	if !result {
		log.Error("Couldn't sign into Discord")
		return false
	}

	return true
}
