package discord

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

type CommandHandler func(s *discordgo.Session, msg *discordgo.Message, cmd, arg string, q bool)

func (srv *Service) AddCommand(name string, fn CommandHandler) {
	srv.Commands[name] = fn
}

func ParseCommand(s string) (cmd, arg string, query bool) {
	s = strings.TrimSpace(s)

	if len(s) == 0 || (s[0] != '/' && s[0] != '!' && s[0] != '?') {
		return "", "", false
	}

	cmdEnd := strings.Index(s, " ")
	if cmdEnd == -1 {
		cmdEnd = len(s)
	}

	return s[1:cmdEnd], strings.TrimSpace(s[cmdEnd:]), s[0] == '?'
}

func (srv *Service) handleMessageCreateWithCommand(s *discordgo.Session, event *discordgo.MessageCreate) {
	cmd, arg, q := ParseCommand(event.Content)
	fn, ok := srv.Commands[cmd]
	if !ok {
		return
	}

	fn(s, event.Message, cmd, arg, q)
}
