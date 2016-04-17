package bot

type Bot struct {
	Services []ServiceRef
}

func New() Bot {
	return Bot{}
}

func (bot *Bot) Run(brain *Brain, stop <-chan interface{}) error {
	for _, service := range bot.Services {
		store, err := brain.Get(TypeService, service.Load)
		if err != nil {
			return err
		}
		service.Impl.Start(store)
		defer service.Impl.Stop(store)
	}

	// Wait for a stop signal
	<-stop

	return nil
}
