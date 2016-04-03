package bot

type Bot struct {
	Services []ServiceRef
}

func New() Bot {
	return Bot{}
}

func (bot *Bot) Run() {
	for i := range bot.Services {
		bot.Services[i].Impl.Start()
	}
}
