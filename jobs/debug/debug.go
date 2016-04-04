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
	session := srv.Session
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		channel, err := session.State.Channel(m.ChannelID)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"id": m.ChannelID,
			}).Error("Couldn't get message's channel")
			return
		}
		log.WithFields(log.Fields{
			"text":    m.Content,
			"channel": channel.Name,
		}).Info("Message")
	})
}
