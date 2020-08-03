package distribution

import (
	common "../common"
	infra "../infra"
	"fmt"
)

type Requestor struct {
	crh               *infra.ClientRequestHandler
	transportType     string
	lastAddress       string
}

func NewRequestor(transportType string) *Requestor {
	return &Requestor{
		transportType:     transportType,
	}
}

func (r *Requestor) Invoke(invocation *common.Invocation) *common.Termination {
	fmt.Println("Requestor: invoking "+invocation.Message.Operation)
	if r.crh == nil || r.lastAddress != invocation.Addr {
		r.crh = infra.NewCRH(invocation.Addr)
	}
	packet := common.NewRequestPacket(*invocation.Message)
	marshalledMessage, err := common.Marshall(*packet)
	if err != nil {return marshallingError(err)}
	if r.transportType == "tcp" {
		return r.SendReceiveTcp(marshalledMessage)
	} else {
		return r.SendReceiveUdp(marshalledMessage)
	}
}

func marshallingError(err error) *common.Termination {
	return &common.Termination{
		Result: fmt.Sprintf("%s: %s", "marshalling error: ", err),
	}
}

func (r *Requestor) SendReceiveTcp(message []byte) *common.Termination {
	r.crh.SendTcp(message)
	data := r.crh.ReceiveTcp()
	return unpackTermination(data)
}

func unpackTermination(data []byte) *common.Termination {
	packet := &common.Packet{}
	err := common.Unmarshall(data, packet)
	if err != nil {return marshallingError(err)}
	ter := string(packet.Body)
	return &common.Termination{Result: ter}
}

func (r *Requestor) SendReceiveUdp(message []byte) *common.Termination {
	r.crh.SendUdp(message)
	data := r.crh.ReceiveUdp()
	return unpackTermination(data)
}