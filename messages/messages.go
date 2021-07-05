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

func DefaultPortConfig() *PortConfig {
	return &PortConfig{
		BaudRate: 115200,
		DataBits: 8,
		StopBits: 0,
		Parity:   "None",
	}
}

type ReconfigurePortRequest struct {
	Config *PortConfig
}

type ReconfigurePortResponse struct {
	Config *PortConfig
}

type ReconnectRequest struct {
	Config *PortConfig
	Port   string
}

type ReconnectResponse struct {
	Config *PortConfig
	Port   string
}
