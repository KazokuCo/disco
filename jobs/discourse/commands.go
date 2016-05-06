package discourse

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"html"
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

	keywords := strings.Split(arg, " ")
	topicLines := []string{}

topicLoop:
	for _, t := range res.Topics {
		lowTitle := strings.ToLower(t.Title)
		for _, keyword := range keywords {
			if !strings.Contains(lowTitle, keyword) {
				continue topicLoop
			}
		}

		title := html.UnescapeString(t.FancyTitle)

		url := fmt.Sprintf("%s/t/%d", j.URL, t.ID)
		line := fmt.Sprintf("%s - <%s>", title, url)
		topicLines = append(topicLines, line)
		log.WithFields(log.Fields{"title": title, "url": url}).Debug("Found topic")
	}

	if len(topicLines) > 3 {
		extra := len(topicLines) - 3
		topicLines = topicLines[:3]
		topicLines = append(topicLines, fmt.Sprintf("+ %d more", extra))
	}
	text := strings.Join(topicLines, "\n")
	s.ChannelMessageSend(msg.ChannelID, text)
}
