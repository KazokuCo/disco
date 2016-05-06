package discourse

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func (j *Job) CommandQueryTopics(s *discordgo.Session, msg *discordgo.Message, cmd, arg string) {
	q := strings.ToLower(arg)
	log.WithField("q", q).Info("Searching for topics")

	res, err := j.Search(arg)
	if err != nil {
		log.WithError(err).Error("Couldn't search forum")
		return
	}

	for _, t := range res.Topics {
		if !strings.Contains(strings.ToLower(t.Title), q) {
			continue
		}

		url := fmt.Sprintf("%s/t/%d", j.URL, t.ID)
		line := fmt.Sprintf("%s - <%s>", t.Title, url)
		s.ChannelMessageSend(msg.ChannelID, line)
		log.WithFields(log.Fields{"title": t.Title, "url": url}).Debug("Found topic")
	}
}
