package discord

import (
	"regexp"
	"testing"
)

func TestListenerMatch(t *testing.T) {
	m := Listener{
		Regex: regexp.MustCompile(`\$([\d]+)\.([\d]+)`),
	}.Match("$10.50")
	if len(m) != 1 || len(m[0]) != 3 {
		t.Fail()
	}
}
