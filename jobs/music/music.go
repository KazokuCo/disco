package music

import (
	"github.com/kazokuco/disco/bot"
)

func init() {
	bot.RegisterJob("music", func() bot.Job { return New() })
}

type Job struct {
	Channel string
	Discord Discord
}

func New() *Job {
	return &Job{}
}
