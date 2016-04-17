package discord

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

type CommandHandler func(s *discordgo.Session, msg *discordgo.Message, offset int)

func (srv *Service) AddCommand(name string, fn CommandHandler) {
	srv.Commands[name] = fn
}

func parseCommand(s string) (cmd string, offset int, ok bool) {
	if len(s) == 0 || (s[0] != '/' && s[0] != '!') {
		return "", 0, false
	}

	cmdEnd := strings.Index(s, " ")
	offset = cmdEnd + 1
	if cmdEnd == -1 {
		cmdEnd = len(s)
		offset = cmdEnd
	}
	return s[1:cmdEnd], offset, true
}

func (srv *Service) handleMessageCreateWithCommand(s *discordgo.Session, event *discordgo.MessageCreate) {
	cmd, offset, ok := parseCommand(event.Content)
	if !ok {
		return
	}

	fn, ok := srv.Commands[cmd]
	if !ok {
		return
	}

	fn(s, event.Message, offset)
}
