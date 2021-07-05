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
	case *messages.PortsRequest:
		ser.loadPorts()
		index, name := ser.selected()
		*ser.channel <- messages.PortsResponse{
			Ports:     ser.ports,
			OpenIndex: index,
			OpenName:  name,
		}
		return true
	case *messages.ReconfigurePortRequest:
		config := msg.(*messages.ReconfigurePortRequest).Config
		*ser.channel <- messages.ReconfigurePortResponse{Config: config}
		return true
	case *messages.ExitRequest:
		*ser.channel <- messages.ExitResponse{}
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
