package discord

import (
	"bufio"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/kazokuco/disco/util"
	"os"
)

func (srv *Service) Authorize(c *cli.Context) {
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

	// Ask for a client ID
	clientID := os.Getenv("DISCORD_CLIENT_ID")
	if clientID == "" {
		if st.Auth.ClientID != "" {
			fmt.Printf("Press ENTER to keep using: %s\n", st.Auth.ClientID)
		}

		fmt.Printf("Client ID: ")
		if !s.Scan() {
			os.Exit(0)
		}
		text := s.Text()
		if text != "" {
			clientID = text
			st.Auth.ClientID = text
		}
	}

	// Use it to generate an authorization link
	link := fmt.Sprintf("https://discordapp.com/oauth2/authorize?&client_id=%s&scope=bot&permissions=0", clientID)
	fmt.Printf("\n%s\n", link)

	if err = util.StoreBrainToFile(&brain, brainFilename); err != nil {
		log.WithError(err).Fatal("Couldn't store brain")
	}
}
