package verification

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/kazokuco/disco/bot"
	"github.com/kazokuco/disco/services/discord"
	"github.com/mvdan/xurls"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func init() {
	bot.RegisterJob("verification", func() bot.Job { return New() })
}

type Job struct {
	Channel   string
	Grant     string
	Against   string
	Discourse struct {
		URL     string
		TopicID int `yaml:"topic_id"`
	}
	Lines struct {
		Success       string
		NameNotInPost string `yaml:"name_not_in_post"`
		Error         string
	}
}

func New() *Job {
	job := Job{}
	job.Lines.Success = "[SUCCESS]"
	job.Lines.NameNotInPost = "[NAME NOT IN POST]"
	job.Lines.Error = "[ERROR]"
	return &job
}

func (j *Job) DiscordInit(srv *discord.Service) {
	session := srv.Session
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// This doesn't get me a GuildID for some reason?
		// channel, err := session.State.Channel(m.ChannelID)
		channel, err := session.Channel(m.ChannelID)
		if err != nil {
			log.WithError(err).WithField("id", m.ChannelID).Error("Couldn't get message's channel")
			return
		}

		guild, err := session.State.Guild(channel.GuildID)
		if err != nil {
			log.WithError(err).WithField("id", channel.GuildID).Error("Couldn't get message's guild")
			return
		}

		if channel.Name != j.Channel {
			return
		}

		grant := &discordgo.Role{}
		for i := range guild.Roles {
			if guild.Roles[i].Name == j.Grant {
				grant = guild.Roles[i]
			}
		}
		if grant.ID == "" {
			log.WithError(err).WithField("name", j.Grant).Error("Couldn't get role to grant")
			return
		}
		log.WithField("grant", grant).Info("Grant")

		urls := xurls.Strict.FindAllString(m.Content, -1)
		for i := range urls {
			if !strings.HasPrefix(urls[i], j.Discourse.URL+"/t/") {
				continue
			}

			u, err := url.Parse(urls[i])
			if err != nil {
				log.WithError(err).WithField("url", u).Warn("Couldn't parse URL")
				continue
			}

			pathParts := strings.Split(u.Path, "/")
			postID, err := strconv.Atoi(pathParts[len(pathParts)-1])
			if err != nil {
				log.WithError(err).Warn("Couldn't get post ID")
			}

			u.RawQuery = ""
			jsonURL := u.String() + ".json"
			res, err := http.Get(jsonURL)
			if err != nil || (res.StatusCode != 200 && res.StatusCode != 404) {
				log.WithError(err).WithField("status", res.StatusCode).Warn("Couldn't fetch JSON")
				srv.Reply(m.Message, j.Lines.Error)
				continue
			}
			defer res.Body.Close()

			body, _ := ioutil.ReadAll(res.Body)
			var data struct {
				ID         int `json:"id"`
				PostStream struct {
					Posts []struct {
						PostNumber int    `json:"post_number"`
						Cooked     string `json:"cooked"`
					} `json:"posts"`
				} `json:"post_stream"`
			}
			err = json.Unmarshal(body, &data)
			if err != nil {
				log.WithError(err).Warn("Couldn't parse JSON")
				srv.Reply(m.Message, j.Lines.Error)
				continue
			}

			for i := range data.PostStream.Posts {
				post := data.PostStream.Posts[i]
				if post.PostNumber == postID {
					if !strings.Contains(post.Cooked, m.Author.Username) {
						srv.Reply(m.Message, j.Lines.NameNotInPost)
						break
					}

					member, err := srv.Session.State.Member(channel.GuildID, m.Author.ID)
					log.WithField("member", member).Info("Member")
					if err != nil {
						log.WithError(err).Warn("Couldn't get member info")
						break
					}
					roles := append(member.Roles, grant.ID)
					err = srv.Session.GuildMemberEdit(channel.GuildID, member.User.ID, roles)
					if err != nil {
						log.WithError(err).WithFields(log.Fields{
							"gid":   member.GuildID,
							"uid":   member.User.ID,
							"roles": roles,
						}).Error("Couldn't grant role")
						srv.Reply(m.Message, j.Lines.Error)
						break
					}

					srv.Reply(m.Message, j.Lines.Success)
				}
			}
		}
	})
}
