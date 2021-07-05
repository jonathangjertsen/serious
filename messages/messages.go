package messages

const (
	ParityNone           = "None"
	ParityEven           = "Even"
	ParityOdd            = "Odd"
	ParityAlways0        = "Always 0"
	ParityAlways1        = "Always 1"
	StopBitsOne          = "1"
	StopBitsOnePointFive = "1.5"
	StopBitsTwo          = "2"
)

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
	StopBits string
	Parity   string
}

func DefaultPortConfig() *PortConfig {
	return &PortConfig{
		BaudRate: 115200,
		DataBits: 8,
		StopBits: StopBitsOne,
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
	Error  error
}

type ReadRequest struct {
	Buffer []byte
	Size   int
}

type ReadResponse struct {
	Buffer []byte
	Size   int
	Error  error
}
