package hw

import (
	messages "github.com/jonathangjertsen/serious/messages"
	"go.bug.st/serial"
)

type Serial struct {
	ports   []string
	channel *chan messages.Message
}

func NewSerial(hwChannel *chan messages.Message) (*Serial, error) {
	ser := Serial{channel: hwChannel}
	err := ser.loadPorts()
	if err != nil {
		return nil, err
	}
	return &ser, nil
}

func (ser *Serial) Run() {
	for {
		message := <-*ser.channel
		keepgoing := ser.handle(message)
		if !keepgoing {
			break
		}
	}
}

func (ser *Serial) handle(msg messages.Message) bool {
	switch msg.(type) {
	case *messages.RequestPorts:
		ser.loadPorts()
		index, name := ser.selected()
		*ser.channel <- messages.RequestPortsResponse{
			Ports:     ser.ports,
			OpenIndex: index,
			OpenName:  name,
		}
		return true
	case *messages.RequestReconfigurePort:
		config := msg.(*messages.RequestReconfigurePort).Config
		*ser.channel <- messages.RequestReconfigurePortResponse{Config: config}
		return true
	case *messages.RequestExit:
		*ser.channel <- messages.RequestExitResponse{}
		return false
	default:
		*ser.channel <- messages.Unexpected{
			Original: msg,
		}
		return false
	}
}

func (ser *Serial) loadPorts() error {
	ports, err := serial.GetPortsList()
	if err != nil {
		return err
	}
	ser.ports = ports
	return nil
}

func (ser *Serial) selected() (int, *string) {
	if len(ser.ports) > 0 {
		return 0, &ser.ports[0]
	} else {
		return 0, nil
	}
}
