package discord

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kazokuco/disco/bot"
)

func init() {
	bot.RegisterService("discord", func() bot.Service { return New() })
}

type Discord struct {
	Jobs []bot.JobRef
}

func New() *Discord {
	return &Discord{}
}

func (dis *Discord) Start() {
	log.Info("Discord: Starting...")
}
