package discord

import (
	"github.com/bwmarrin/discordgo"
	"regexp"
)

type ListenerAction func(s *discordgo.Session, msg *discordgo.Message, matches [][]string)

type Listener struct {
	Regexp *regexp.Regexp
	Action ListenerAction
}

func (l Listener) Match(s string) [][]string {
	// Don't interfere with custom commands; do react to /me "commands"
	if (s[0] == '/' || s[0] == '!') && s[:3] != "/me" {
		return nil
	}

	return l.Regexp.FindAllStringSubmatch(s, -1)
}

func (srv *Service) AddListener(l Listener) {
	srv.Listeners = append(srv.Listeners, l)
}

func (srv *Service) handleMessageCreateWithListener(s *discordgo.Session, event *discordgo.MessageCreate) {
	for _, l := range srv.Listeners {
		m := l.Match(event.Content)
		if m != nil {
			l.Action(s, event.Message, m)
		}
	}
}
