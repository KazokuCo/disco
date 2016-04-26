package verification

import (
	"encoding/json"
	"errors"
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
	Discourse struct {
		RawURL     string   `yaml:"url"`
		TopicID    int      `yaml:"topic_id"`
		TrustLevel int      `yaml:"trust_level"`
		URL        *url.URL `yaml:"-"`
	}
	Lines struct {
		Success       string
		NameNotInPost string `yaml:"name_not_in_post"`
		LevelTooLow   string `yaml:"level_too_low"`
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

func (j *Job) Init() (err error) {
	// Parse the forum base URL
	if j.Discourse.RawURL != "" {
		j.Discourse.URL, err = url.Parse(j.Discourse.RawURL)
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *Job) DiscordInit(srv *discord.Service) {
	if err := j.Init(); err != nil {
		log.WithError(err).Error("Verification Init failed")
		return
	}

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
			log.WithField("channel", j.Channel).Debug("Verification: Wrong channel")
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

		urls, err := ParseLinks(m.Content, j.Discourse.URL)
		if err != nil {
			log.WithError(err).Error("Couldn't parse message URLs")
			return
		}
		for i := range urls {
			u := urls[i]
			log.WithField("url", u).Debug("Verification: URL found")

			// Ignore links not to topics
			if !strings.HasPrefix(u.Path, "/t/") {
				continue
			}

			// Parse the topic ID from the URL
			topicID, postID, err := ParseDiscourseTopicURL(u)
			if err != nil {
				log.WithError(err).Warn("Failed to parse post URL")
				continue
			}

			// Ignore posts in the wrong topic
			if topicID != j.Discourse.TopicID {
				continue
			}

			// Get the JSON for the topic
			res, err := GetJSON(u.String())
			defer res.Body.Close()
			if err != nil || (res.StatusCode != 200 && res.StatusCode != 404) {
				log.WithError(err).WithField("status", res.StatusCode).Warn("Couldn't fetch JSON")
				srv.Reply(m.Message, j.Lines.Error)
				continue
			}
			body, _ := ioutil.ReadAll(res.Body)

			// Parse it into a subset of the topic JSON structure
			var data struct {
				ID         int `json:"id"`
				PostStream struct {
					Posts []struct {
						PostNumber int    `json:"post_number"`
						Cooked     string `json:"cooked"`
						TrustLevel int    `json:"trust_level"`
					} `json:"posts"`
				} `json:"post_stream"`
			}
			err = json.Unmarshal(body, &data)
			if err != nil {
				log.WithError(err).Warn("Couldn't parse JSON")
				srv.Reply(m.Message, j.Lines.Error)
				continue
			}

			// Double-check that this is the right topic
			if data.ID != j.Discourse.TopicID {
				continue
			}

			// Loop through posts until we find the one linked one
			for i := range data.PostStream.Posts {
				post := data.PostStream.Posts[i]
				if post.PostNumber != postID {
					continue
				}

				// Make sure the correct Discord username is mentioned
				if !strings.Contains(post.Cooked, m.Author.Username) {
					srv.Reply(m.Message, j.Lines.NameNotInPost)
					break
				}

				// If requested, verify the user's trust level
				if j.Discourse.TrustLevel > post.TrustLevel {
					log.WithFields(log.Fields{
						"have": post.TrustLevel,
						"need": j.Discourse.TrustLevel,
					}).Debug("Verification: Trust Level")
					srv.Reply(m.Message, j.Lines.LevelTooLow)
					break
				}

				member, err := srv.Session.State.Member(channel.GuildID, m.Author.ID)
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
	})
}

func ParseLinks(text string, baseURL *url.URL) (urls []*url.URL, err error) {
	// Find all URLs in the post
	raw := xurls.Strict.FindAllString(text, -1)
	for i := range raw {
		// On the off chance that the URL does not parse, skip it
		u, err := url.Parse(raw[i])
		if err != nil {
			continue
		}
		// Only add URLs matching the configured host
		if u.Host == baseURL.Host {
			urls = append(urls, u)
		}
	}
	return urls, nil
}

func ParseDiscourseTopicURL(u *url.URL) (topic, post int, err error) {
	parts := strings.Split(u.Path, "/")
	if len(parts) < 2 {
		return 0, 0, errors.New("Not enough parts in the URL")
	}
	topic, err = strconv.Atoi(parts[len(parts)-2])
	if err != nil {
		return topic, post, err
	}
	post, err = strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return topic, post, err
	}
	return topic, post, nil
}

func GetJSON(u string) (res *http.Response, err error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return res, err
	}
	req.Header.Add("Accept", "application/json")
	return client.Do(req)
}
