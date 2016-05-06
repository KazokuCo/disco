package discourse

import (
	"encoding/json"
	"github.com/kazokuco/disco/bot"
	"github.com/kazokuco/disco/services/discord"
	"io/ioutil"
	"net/http"
	"net/url"
)

func init() {
	bot.RegisterJob("discourse", func() bot.Job { return New() })
}

type Job struct {
	URL string
}

func New() *Job {
	return &Job{}
}

func (j *Job) DiscordInit(srv *discord.Service) {
	srv.AddCommand("?topic", j.CommandQueryTopics)
}

func (j *Job) Search(q string) (env SearchEnvelope, err error) {
	res, err := http.Get(j.URL + "/search.json?q=" + url.QueryEscape(q))
	if err != nil {
		return env, err
	}
	defer res.Body.Close()

	data, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(data, &env); err != nil {
		return env, err
	}

	return env, nil
}
