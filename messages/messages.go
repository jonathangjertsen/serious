package messages

type Message interface {
}

type Unexpected struct {
	Original Message
}

type RequestPorts struct {
}

type RequestPortsResponse struct {
	Ports     []string
	OpenName  *string
	OpenIndex int
}

type RequestExit struct {
}

type RequestExitResponse struct {
}

type PortConfig struct {
	BaudRate int
	DataBits int
	StopBits int
	Parity   string
}

type RequestReconfigurePort struct {
	Config *PortConfig
}

type RequestReconfigurePortResponse struct {
	Config *PortConfig
}
