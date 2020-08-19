package distribution

import (
	common "../common"
	"fmt"
)

type ClientProxy struct {
	requestor   *Requestor
	srvAddress  string
	id    string
}

func NewClientProxy(aor *common.AOR) *ClientProxy {
	return &ClientProxy{
		requestor:   NewRequestor(aor.Protocol),
		srvAddress:  aor.Address,
		id:          aor.ObjectId,
	}
}

func (cp *ClientProxy) ChangeLane(newLane string) {
	cp.invokeCommand("CHANGE", newLane)
}

func (cp *ClientProxy) BroadcastEvent(lane string) {
	cp.invokeCommand("BroadcastEvent", lane)
}

func (cp *ClientProxy) RegisterOnLane(lane string) {
	cp.invokeCommand("Register", lane)
}

func (cp *ClientProxy) invokeCommand(op string, lane string) {
	fmt.Println("ClientProxy: invoking "+op)
	var message = &common.Message{
		Operation: op,
		Topic:     lane,
	}
	var invocation = &common.Invocation{
		Addr: cp.srvAddress,
		Message: message,
	}
	cp.requestor.Invoke(invocation)
}


