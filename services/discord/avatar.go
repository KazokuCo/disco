package discord

import (
	"encoding/base64"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/kazokuco/disco/bot"
	"io/ioutil"
	"os"
)

func (srv *Service) UpdateAvatar(store bot.Store, filename string) {
	st := store.(*Store)
	token := st.Auth.Token
	if token == "" {
		token = os.Getenv("DISCORD_TOKEN")
	}
	session, err := discordgo.New(token)
	if err != nil {
		log.WithError(err).Fatal("Couldn't connect to Discord")
	}

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

	result := <-resultCh
	if !result {
		log.Fatal("Couldn't sign into Discord")
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithError(err).Fatal("Couldn't read avatar")
	}
	b64 := base64.StdEncoding.EncodeToString(data)

	avatarData := fmt.Sprintf("data:image/jpeg;base64,%s", b64)
	_, err = session.UserUpdate(session.State.User.Email, "", session.State.User.Username, avatarData, "")
	if err != nil {
		log.WithError(err).Fatal("Couldn't update avatar")
	}
}
