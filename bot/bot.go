package bot

type Bot struct {
	Services []ServiceRef
}

func New() Bot {
	return Bot{}
}

func (bot *Bot) Run(brain Brain) error {
	for i := range bot.Services {
		service := bot.Services[i]
		store, err := brain.Get(TypeService, service.Load)
		if err != nil {
			return err
		}
		service.Impl.Start(store)
	}

	return nil
}
