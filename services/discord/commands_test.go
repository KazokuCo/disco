package discord

import (
	"testing"
)

func TestParseCommandBlank(t *testing.T) {
	cmd, arg, q := ParseCommand("")
	if cmd != "" || arg != "" || q {
		t.Fail()
	}
}

func TestParseCommandNotACommand(t *testing.T) {
	cmd, arg, q := ParseCommand("cmd arg")
	if cmd != "" || arg != "" || q {
		t.Fail()
	}
}

func TestParseCommandSingle(t *testing.T) {
	cmd, arg, q := ParseCommand("!cmd")
	if cmd != "cmd" || arg != "" || q {
		t.Fail()
	}
}

func TestParseCommandArg(t *testing.T) {
	cmd, arg, q := ParseCommand("!cmd arg")
	if cmd != "cmd" || arg != "arg" || q {
		t.Fail()
	}
}

func TestParseCommandRedundantSpaces(t *testing.T) {
	cmd, arg, q := ParseCommand("!cmd   arg   ")
	if cmd != "cmd" || arg != "arg" || q {
		t.Fail()
	}
}

func TestParseCommandQuery(t *testing.T) {
	cmd, arg, q := ParseCommand("?cmd query")
	if cmd != "cmd" || arg != "query" || !q {
		t.Fail()
	}
}
