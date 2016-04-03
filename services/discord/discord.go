package discord

import (
	"bufio"
	log "github.com/Sirupsen/logrus"
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
	}
}

func New() *Service {
	return &Service{}
}

func (srv *Service) Authorize(store bot.Store) bool {
	st := store.(*Store)
	s := bufio.NewScanner(os.Stdin)

	print("username: ")
	if !s.Scan() {
		return false
	}
	st.Auth.Username = s.Text()

	print("password: ")
	if !s.Scan() {
		return false
	}
	st.Auth.Password = s.Text()

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
