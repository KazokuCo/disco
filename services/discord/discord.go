package discord

import (
	"bufio"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/kazokuco/disco/bot"
	"os"
	"strings"
)

func init() {
	bot.RegisterService("discord", func() bot.Service { return New() })
}

type Job interface {
	DiscordInit(srv *Service)
}

type Service struct {
	Game string
	Jobs []bot.JobRef

	Session *discordgo.Session `yaml:"-"`
}

type Store struct {
	Auth struct {
		Token string
	}
}

func New() *Service {
	return &Service{}
}

func (srv *Service) Login(store bot.Store) bool {
	st := store.(*Store)
	s := bufio.NewScanner(os.Stdin)

	print("Token: ")
	if !s.Scan() {
		return false
	}
	st.Auth.Token = s.Text()

	_, err := discordgo.New(st.Auth.Token)
	if err != nil {
		log.WithError(err).Error("Couldn't sign into Discord")
		return false
	}

	return true
}

func (srv *Service) Store() bot.Store {
	return &Store{}
}

func (srv *Service) Start(store bot.Store) {
	st := store.(*Store)

	token := st.Auth.Token
	if token == "" {
		token = os.Getenv("DISCORD_TOKEN")
	}
	if token == "" {
		log.Fatal("No Discord token given!")
	}
	log.Info("Discord: Starting...")

	session, err := discordgo.New(token)
	if err != nil {
		log.WithError(err).Fatal("Couldn't connect to Discord")
	}
	srv.Session = session

	for i := range srv.Jobs {
		ref := srv.Jobs[i]
		job, ok := ref.Impl.(Job)
		if !ok {
			log.WithField("job", ref.Load).Fatal("Job does not support Discord")
		}
		job.DiscordInit(srv)
	}

	if err = srv.Session.Open(); err != nil {
		log.WithError(err).Fatal("Discord: Failed to open connection!")
		return
	}

	if err = srv.Session.UpdateStatus(0, srv.Game); err != nil {
		log.WithError(err).Warn("Discord: Failed to update status")
	}
}

func (srv *Service) MentionsMe(m *discordgo.Message) bool {
	mention := fmt.Sprintf("<@%s>", srv.Session.State.User.ID)
	return strings.Contains(m.Content, mention)
}

func (srv *Service) Reply(m *discordgo.Message, text string) {
	srv.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> %s", m.Author.ID, text))
}
