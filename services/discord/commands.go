package discord

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

type CommandHandler func(s *discordgo.Session, msg *discordgo.Message, cmd, arg string)

func (srv *Service) AddCommand(name string, fn CommandHandler) {
	srv.Commands[name] = fn
}

func ParseCommand(s string) (cmd, arg string) {
	s = strings.TrimSpace(s)

	if len(s) == 0 || (s[0] != '!' && s[0] != '?') {
		return "", ""
	}

	cmdEnd := strings.Index(s, " ")
	if cmdEnd == -1 {
		cmdEnd = len(s)
	}

	return s[:cmdEnd], strings.TrimSpace(s[cmdEnd:])
}

func (srv *Service) handleMessageCreateWithCommand(s *discordgo.Session, event *discordgo.MessageCreate) {
	cmd, arg := ParseCommand(event.Content)
	fn, ok := srv.Commands[cmd]
	if !ok {
		return
	}

	fn(s, event.Message, cmd, arg)
}
