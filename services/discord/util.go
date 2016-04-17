package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func (srv *Service) MentionsMe(m *discordgo.Message) bool {
	mention := fmt.Sprintf("<@%s>", srv.Session.State.User.ID)
	return strings.Contains(m.Content, mention)
}

func (srv *Service) Reply(m *discordgo.Message, text string) {
	srv.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> %s", m.Author.ID, text))
}
