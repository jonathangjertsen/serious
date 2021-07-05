package hw

import (
	"go.bug.st/serial"
)

type Serial struct {
	ports []string
}

func NewSerial() (*Serial, error) {
	ser := Serial{}
	err := ser.LoadPorts()
	if err != nil {
		return nil, err
	}
	return &ser, nil
}

func (ser *Serial) LoadPorts() error {
	ports, err := serial.GetPortsList()
	if err != nil {
		return err
	}
	ser.ports = ports
	return nil
}

func (ser *Serial) GetPorts() []string {
	return ser.ports
}

func (ser *Serial) Selected() (int, *string) {
	if len(ser.ports) > 0 {
		return 0, &ser.ports[0]
	} else {
		return 0, nil
	}
}
