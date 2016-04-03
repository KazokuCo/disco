package discord

import (
	"bufio"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/kazokuco/disco/bot"
	"os"
)

func init() {
	bot.RegisterService("discord", func() bot.Service { return New() })
}

type Service struct {
	Jobs []bot.JobRef
}

type Store struct {
	Auth struct {
		Username string
		Password string
		Token    string
	}
}

func New() *Service {
	return &Service{}
}

func (srv *Service) Login(store bot.Store) bool {
	st := store.(*Store)
	s := bufio.NewScanner(os.Stdin)

	print("Username: ")
	if !s.Scan() {
		return false
	}
	st.Auth.Username = s.Text()

	print("Password: ")
	if !s.Scan() {
		return false
	}
	st.Auth.Password = s.Text()

	session, err := discordgo.New(st.Auth.Username, st.Auth.Password)
	if err != nil {
		log.WithError(err).Error("Couldn't sign into Discord")
		return false
	}

	st.Auth.Token = session.Token

	return true
}

func (srv *Service) Store() bot.Store {
	return &Store{}
}

func (srv *Service) Start(store bot.Store) {
	st := store.(*Store)
	log.Info("Discord: Starting...")
	log.WithFields(log.Fields{
		"username": st.Auth.Username,
	}).Info("Auth")
}
