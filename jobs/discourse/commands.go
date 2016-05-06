package discourse

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"net/html"
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

	topicLines := []string{}
	for _, t := range res.Topics {
		if !strings.Contains(strings.ToLower(t.Title), q) {
			continue
		}

		url := fmt.Sprintf("%s/t/%d", j.URL, t.ID)
		line := fmt.Sprintf("%s - <%s>", t.Title, url)
		topicLines = append(topicLines, line)
		log.WithFields(log.Fields{"title": t.Title, "url": url}).Debug("Found topic")
	}

	text := strings.Join(topicLines, "\n")
	s.ChannelMessageSend(msg.ChannelID, text)
}