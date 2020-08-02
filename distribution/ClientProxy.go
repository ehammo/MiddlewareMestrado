package distribution

import (
	common "../common"
	"fmt"
)

type ClientProxy struct {
	requestor   *Requestor
	srvAddress  string
	typeName    string
}

func NewClientProxy(srvAddress string,
	                transportType string) *ClientProxy {
	return &ClientProxy{
		requestor:   NewRequestor(transportType),
		srvAddress:  srvAddress,
	}
}

func (cp *ClientProxy) ChangeLane(newLane string) string {
	return cp.invokeCommand("CHANGE", newLane)
}

func (cp *ClientProxy) BroadcastEvent(lane string) string {
	return cp.invokeCommand("BREAK", lane)
}

func (cp *ClientProxy) RegisterOnLane(lane string) string {
	return cp.invokeCommand("REGISTER", lane)
}

func (cp *ClientProxy) invokeCommand(op string, lane string) string {
	fmt.Println("ClientProxy: invoking "+op)
	var message = &common.Message{
		Operation: op,
		Topic:     lane,
	}
	var invocation = &common.Invocation{
		Addr: cp.srvAddress,
		Message: message,
	}
	response := cp.requestor.Invoke(invocation)
	return response.Result
}


