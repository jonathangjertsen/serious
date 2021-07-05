package messages

type Message interface {
}

type Unexpected struct {
	Original Message
}

type PortsRequest struct {
}

type PortsResponse struct {
	Ports     []string
	OpenName  string
	OpenIndex int
}

type ExitRequest struct {
}

type ExitResponse struct {
}

type PortConfig struct {
	BaudRate int
	DataBits int
	StopBits int
	Parity   string
}

type ReconfigurePortRequest struct {
	Config *PortConfig
}

type ReconfigurePortResponse struct {
	Config *PortConfig
}
