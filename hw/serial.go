package hw

import (
	messages "github.com/jonathangjertsen/serious/messages"
	"go.bug.st/serial"
)

type Serial struct {
	ports     []string
	openIndex int
	openName  string
	channel   *chan messages.Message
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
		if !ser.handle(message) {
			break
		}
	}
}

func (ser *Serial) handle(msg messages.Message) bool {
	switch msg.(type) {
	case *messages.PortsRequest:
		ser.loadPorts()
		*ser.channel <- messages.PortsResponse{
			Ports:     ser.ports,
			OpenIndex: ser.openIndex,
			OpenName:  ser.openName,
		}
		return true
	case *messages.ReconfigurePortRequest:
		config := msg.(*messages.ReconfigurePortRequest).Config
		*ser.channel <- messages.ReconfigurePortResponse{Config: config}
		return true
	case *messages.ReconnectRequest:
		request := msg.(*messages.ReconnectRequest)
		for index, port := range ser.ports {
			if port == request.Port {
				ser.openIndex = index
				ser.openName = port
				break
			}
		}
		*ser.channel <- messages.ReconnectResponse{Config: request.Config, Port: ser.openName}
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
