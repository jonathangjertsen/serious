package hw

import (
	"fmt"
	messages "github.com/jonathangjertsen/serious/messages"
	"go.bug.st/serial"
)

var ParityMapping = map[string]serial.Parity{
	messages.ParityNone:    serial.NoParity,
	messages.ParityEven:    serial.EvenParity,
	messages.ParityOdd:     serial.OddParity,
	messages.ParityAlways1: serial.MarkParity,
	messages.ParityAlways0: serial.SpaceParity,
}

var StopBitsMapping = map[string]serial.StopBits{
	messages.StopBitsOne:          serial.OneStopBit,
	messages.StopBitsOnePointFive: serial.OnePointFiveStopBits,
	messages.StopBitsTwo:          serial.TwoStopBits,
}

type Serial struct {
	ports     []string
	openIndex int
	openName  string
	open      serial.Port
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
		var err error
		request := msg.(*messages.ReconnectRequest)
		for index, port := range ser.ports {
			if port == request.Port {
				err = ser.connect(port, request.Config)
				if err == nil {
					ser.openIndex = index
					ser.openName = port
				}
				break
			}
		}
		*ser.channel <- messages.ReconnectResponse{
			Config: request.Config,
			Port:   ser.openName,
			Error:  err,
		}
		return true
	case *messages.ReadRequest:
		request := msg.(*messages.ReadRequest)
		size, err := ser.read(request.Buffer, request.Size)
		*ser.channel <- messages.ReadResponse{
			Buffer: request.Buffer,
			Size:   size,
			Error:  err,
		}
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

func (ser *Serial) connect(port string, config *messages.PortConfig) error {
	mode := &serial.Mode{
		BaudRate: config.BaudRate,
		DataBits: config.DataBits,
		Parity:   ParityMapping[config.Parity],
		StopBits: StopBitsMapping[config.StopBits],
	}
	portObj, err := serial.Open(port, mode)
	ser.open = portObj
	return err
}

func (ser *Serial) read(buffer []byte, size int) (int, error) {
	if ser.open == nil {
		return 0, fmt.Errorf("Attempted to read without a connected port")
	}
	return ser.open.Read(buffer)
}
