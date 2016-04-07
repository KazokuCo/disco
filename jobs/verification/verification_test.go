package verification

import (
	"net/url"
	"testing"
)

func URL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

func TestParseLinks(t *testing.T) {
	base, _ := url.Parse("https://example.com/")
	text := "Go to https://example.com/memes/ for hot new memes."
	urls, _ := ParseLinks(text, base)
	if len(urls) != 1 {
		t.Fail()
	}
}

func TestParseDiscourseTopicURL(t *testing.T) {
	tid, pid, err := ParseDiscourseTopicURL(URL("http://forum.example.com/t/12/4"))
	if err != nil {
		t.Errorf("Couldn't parse: %s", err)
	}
	if tid != 12 {
		t.Errorf("Incorrect thread: %v != 12", tid)
	}
	if pid != 4 {
		t.Errorf("Incorrect post: %v != 4", pid)
	}
}
