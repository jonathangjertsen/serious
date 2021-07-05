package ui

import (
	messages "github.com/jonathangjertsen/serious/messages"
	"testing"
)

func Test_NewTerminal(t *testing.T) {
	term := NewTerminal()

	if term.channel != nil {
		t.Fatal("term.channel should be nil")
	}

	if focused := term.app.GetFocus(); focused == nil {
		t.Fatal("No element is in focus")
	}
}

func Test_getPortConfig(t *testing.T) {
	term := NewTerminal()
	have, err := term.getPortConfig()
	if err != nil {
		t.Fatalf("getPortConfig() returned %s", err)
	}
	want := messages.DefaultPortConfig()
	if *have != *want {
		t.Fatalf("have: %+v\nwant: %+v", have, want)
	}
}
