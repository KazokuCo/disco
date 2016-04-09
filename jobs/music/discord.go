package music

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/kazokuco/disco/services/discord"
)

type Discord struct {
	Voice *discordgo.VoiceConnection
}

func (j *Job) DiscordInit(srv *discord.Service) {
	srv.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Commands must mention the bot
		if !srv.MentionsMe(m.Message) {
			return
		}

		// Find a request to play an audio file
		if filename := FindPlayFilename(m.Content); filename != "" {
			path := "/Users/uppfinnarn/Desktop/" + filename
			if err := j.DiscordPlayFile(srv, path); err != nil {
				log.WithError(err).WithField("path", path).Error("Discord/Music: Couldn't play file")
			}
		}
	})
}

func (j *Job) DiscordPlayFile(srv *discord.Service, path string) error {
	// srv.Session.ChannelVoiceJoin(gID, cID, mute, deaf)
	return nil
}
