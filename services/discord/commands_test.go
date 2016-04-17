package discord

import (
	"testing"
)

func TestParseCommandSingle(t *testing.T) {
	s := "/command"
	cmd, offset, ok := parseCommand(s)
	if !ok {
		t.FailNow()
	}
	if cmd != "command" {
		t.Error("cmd")
	}
	if s[offset:] != "" {
		t.Error("offset")
	}
}

func TestParseCommandTrailingSpace(t *testing.T) {
	s := "/command "
	cmd, offset, ok := parseCommand(s)
	if !ok {
		t.FailNow()
	}
	if cmd != "command" {
		t.Error("cmd")
	}
	if s[offset:] != "" {
		t.Error("offset")
	}
}

func TestParseCommandArgs(t *testing.T) {
	s := "/command something"
	cmd, offset, ok := parseCommand(s)
	if !ok {
		t.FailNow()
	}
	if cmd != "command" {
		t.Error("cmd")
	}
	if s[offset:] != "something" {
		t.Error("offset")
	}
}
