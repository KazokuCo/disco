package conversion

import (
	// log "github.com/Sirupsen/logrus"
	// "github.com/bwmarrin/discordgo"
	"github.com/kazokuco/disco/bot"
	"github.com/kazokuco/disco/services/discord"
)

func init() {
	bot.RegisterJob("conversion", func() bot.Job { return New() })
}

type Job struct {
	Lines struct {
		Currency      string
		CurrencyMulti string `yaml:"currency_multi"`
	}
}

func New() *Job {
	return &Job{}
}

func (j *Job) DiscordInit(srv *discord.Service) {
	srv.AddListener(discord.Listener{
		Regexp: currencyRegex,
		Action: j.HandleCurrency,
	})
}
