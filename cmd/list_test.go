package cmd

import (
	messages "github.com/jonathangjertsen/serious/messages"
	"testing"
)

func Test_portListString(t *testing.T) {
	resp := messages.PortsResponse{
		Ports:     []string{"COM1", "COM2", "COM3"},
		OpenIndex: 0,
		OpenName:  "COM1",
	}
	have := portListString(&resp)
	want := `COM1 [will be auto-selected]
COM2
COM3
`
	if have != want {
		t.Fatalf("have:\n%s\n\nwant:\n%s", have, want)
	}
}
