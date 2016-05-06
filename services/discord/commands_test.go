package discord

import (
	"testing"
)

func TestParseCommandBlank(t *testing.T) {
	cmd, arg := ParseCommand("")
	if cmd != "" || arg != "" {
		t.Fail()
	}
}

func TestParseCommandNotACommand(t *testing.T) {
	cmd, arg := ParseCommand("cmd arg")
	if cmd != "" || arg != "" {
		t.Fail()
	}
}

func TestParseCommandSingle(t *testing.T) {
	cmd, arg := ParseCommand("!cmd")
	if cmd != "!cmd" || arg != "" {
		t.Fail()
	}
}

func TestParseCommandArg(t *testing.T) {
	cmd, arg := ParseCommand("!cmd arg")
	if cmd != "!cmd" || arg != "arg" {
		t.Fail()
	}
}

func TestParseCommandRedundantSpaces(t *testing.T) {
	cmd, arg := ParseCommand("!cmd   arg   ")
	if cmd != "!cmd" || arg != "arg" {
		t.Fail()
	}
}

func TestParseCommandQuery(t *testing.T) {
	cmd, arg := ParseCommand("?cmd query")
	if cmd != "?cmd" || arg != "query" {
		t.Fail()
	}
}
