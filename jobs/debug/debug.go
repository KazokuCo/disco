package debug

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/kazokuco/disco/bot"
	"github.com/kazokuco/disco/services/discord"
)

func init() {
	bot.RegisterJob("debug", func() bot.Job { return New() })
}

type Job struct{}

func New() *Job {
	return &Job{}
}

func (j *Job) DiscordInit(srv *discord.Service) {
	srv.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		log.WithFields(log.Fields{
			"text": m.Content,
		}).Info("Message")
	})
}
