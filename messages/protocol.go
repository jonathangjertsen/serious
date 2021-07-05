package messages

import (
	"sync"
)

var mut sync.Mutex

func SyncGetPorts(channel *chan Message) PortsResponse {
	mut.Lock()
	*channel <- &PortsRequest{}
	response := <-*channel
	mut.Unlock()
	return response.(PortsResponse)
}

func SyncExit(channel *chan Message) ExitResponse {
	mut.Lock()
	*channel <- &ExitRequest{}
	response := <-*channel
	mut.Unlock()
	return response.(ExitResponse)
}

func SyncReconfigurePort(channel *chan Message, config *PortConfig) ReconfigurePortResponse {
	mut.Lock()
	*channel <- &ReconfigurePortRequest{Config: config}
	response := <-*channel
	mut.Unlock()
	return response.(ReconfigurePortResponse)
}

func SyncReconnectPort(channel *chan Message, port string, config *PortConfig) ReconnectResponse {
	mut.Lock()
	*channel <- &ReconnectRequest{Port: port, Config: config}
	response := <-*channel
	mut.Unlock()
	return response.(ReconnectResponse)
}

func SyncRead(channel *chan Message, buffer []byte, size int) ReadResponse {
	mut.Lock()
	*channel <- &ReadRequest{Buffer: buffer, Size: size}
	response := <-*channel
	mut.Unlock()
	return response.(ReadResponse)
}
