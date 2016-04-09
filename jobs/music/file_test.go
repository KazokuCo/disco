package music

import (
	"testing"
)

func TestFindPlayFilename(t *testing.T) {
	if file := FindPlayFilename("play ahaha.wav"); file != "ahaha.wav" {
		t.Errorf("Incorrect filename; got '%s'", file)
	}
}
