package discord

import (
	"regexp"
	"testing"
)

func TestListenerMatch(t *testing.T) {
	m := Listener{
		Regexp: regexp.MustCompile(`\$([\d]+)\.([\d]+)`),
	}.Match("$10.50")
	if len(m) != 1 || len(m[0]) != 3 {
		t.Fail()
	}
}

func TestListenerMatchCustomCommand(t *testing.T) {
	if (Listener{Regexp: regexp.MustCompile(`\$\d+`)}).Match("/cmd $10") != nil {
		t.Fail()
	}
}

func TestListenerMatchSlashMe(t *testing.T) {
	if (Listener{Regexp: regexp.MustCompile(`\$\d+`)}).Match("/me $10") == nil {
		t.Fail()
	}
}
