package discord

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/codegangsta/cli"
	"github.com/kazokuco/disco/bot"
	"os"
)

func init() {
	bot.RegisterService("discord", func() bot.Service { return New() })
}

type Job interface {
	DiscordInit(srv *Service)
}

type Service struct {
	Game string
	Jobs []bot.JobRef

	Session *discordgo.Session `yaml:"-"`
}

type Store struct {
	Auth struct {
		Token    string // Bot user token, for user auth
		ClientID string // App client ID, for generating login URLs
	}
}

func New() *Service {
	return &Service{}
}

func (srv *Service) Store() bot.Store {
	return &Store{}
}

func (srv *Service) Start(store bot.Store) {
	st := store.(*Store)

	token := st.Auth.Token
	if token == "" {
		token = os.Getenv("DISCORD_TOKEN")
	}
	if token == "" {
		log.Fatal("No Discord token given!")
	}
	log.Info("Discord: Starting...")

	session, err := discordgo.New(token)
	if err != nil {
		log.WithError(err).Fatal("Couldn't connect to Discord")
	}
	srv.Session = session

	for i := range srv.Jobs {
		ref := srv.Jobs[i]
		job, ok := ref.Impl.(Job)
		if !ok {
			log.WithField("job", ref.Load).Fatal("Job does not support Discord")
		}
		job.DiscordInit(srv)
	}

	srv.Session.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		log.WithFields(log.Fields{
			"ver":  event.Version,
			"beat": event.HeartbeatInterval,
		}).Info("Discord: Ready!")

		if err = srv.Session.UpdateStatus(0, srv.Game); err != nil {
			log.WithError(err).Warn("Discord: Failed to update status")
		}
	})
	srv.Session.AddHandler(func(s *discordgo.Session, event *discordgo.Disconnect) {
		log.Warn("Discord: Disconnected!")
	})
	srv.Session.AddHandler(func(s *discordgo.Session, event *discordgo.RateLimit) {
		log.WithFields(log.Fields{
			"bucket": event.Bucket,
			"msg":    event.Message,
			"retry":  event.RetryAfter,
		}).Warn("Discord: Rate limited!")
	})

	if err = srv.Session.Open(); err != nil {
		log.WithError(err).Fatal("Discord: Failed to open connection!")
		return
	}
}

func (srv *Service) Command() cli.Command {
	return cli.Command{
		Name:  "discord",
		Usage: "Discord-specific commands",
		Subcommands: []cli.Command{
			cli.Command{
				Name:      "login",
				ArgsUsage: "bot.yml",
				Action:    func(c *cli.Context) { srv.Login(c) },
			},
			cli.Command{
				Name:      "auth",
				ArgsUsage: "bot.yml",
				Action:    func(c *cli.Context) { srv.Authorize(c) },
			},
		},
	}
}
