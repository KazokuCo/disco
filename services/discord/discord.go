package discord

import (
	"bufio"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/kazokuco/disco/bot"
	"os"
)

func init() {
	bot.RegisterService("discord", func() bot.Service { return New() })
}

type Job interface {
	DiscordInit(srv *Service)
}

type Service struct {
	Jobs []bot.JobRef

	Session *discordgo.Session `yaml:"-"`
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
	log.WithField("username", st.Auth.Username).Info("Discord: Starting...")

	// Try to connect using an auth token first, otherwise fall back to username + password; in
	// case the token has expired and been re-issued, always store the latest one
	session, err := discordgo.New(st.Auth.Username, st.Auth.Password, st.Auth.Token)
	if err != nil {
		log.WithError(err).Fatal("Couldn't connect to Discord")
	}
	srv.Session = session
	st.Auth.Token = srv.Session.Token

	for i := range srv.Jobs {
		ref := srv.Jobs[i]
		job, ok := ref.Impl.(Job)
		if !ok {
			log.WithField("job", ref.Load).Fatal("Job does not support Discord")
		}
		job.DiscordInit(srv)
	}

	err = srv.Session.Open()
}

func (srv *Service) Reply(m *discordgo.Message, text string) {
	srv.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> %s", m.Author.ID, text))
}
