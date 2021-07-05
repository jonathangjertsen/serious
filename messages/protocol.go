package messages

func SyncGetPorts(channel *chan Message) RequestPortsResponse {
	*channel <- &RequestPorts{}
	response := <-*channel
	return response.(RequestPortsResponse)
}

func SyncExit(channel *chan Message) RequestExitResponse {
	*channel <- &RequestExit{}
	response := <-*channel
	return response.(RequestExitResponse)
}

func SyncReconfigurePort(channel *chan Message, config *PortConfig) RequestReconfigurePortResponse {
	*channel <- &RequestReconfigurePort{Config: config}
	response := <-*channel
	return response.(RequestReconfigurePortResponse)
}
