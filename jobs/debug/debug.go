package debug

import (
	"github.com/kazokuco/disco/bot"
)

func init() {
	bot.RegisterJob("debug", func() bot.Job { return New() })
}

type Job struct{}

func New() *Job {
	return &Job{}
}
