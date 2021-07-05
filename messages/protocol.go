package messages

func SyncGetPorts(channel *chan Message) PortsResponse {
	*channel <- &PortsRequest{}
	response := <-*channel
	return response.(PortsResponse)
}

func SyncExit(channel *chan Message) ExitResponse {
	*channel <- &ExitRequest{}
	response := <-*channel
	return response.(ExitResponse)
}

func SyncReconfigurePort(channel *chan Message, config *PortConfig) ReconfigurePortResponse {
	*channel <- &ReconfigurePortRequest{Config: config}
	response := <-*channel
	return response.(ReconfigurePortResponse)
}

func SyncReconnectPort(channel *chan Message, port string, config *PortConfig) ReconnectResponse {
	*channel <- &ReconnectRequest{Port: port, Config: config}
	response := <-*channel
	return response.(ReconnectResponse)
}
