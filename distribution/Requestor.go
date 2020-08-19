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
	message := *invocation.Message
	fmt.Println("Requestor: invoking "+message.Operation)
	if r.crh == nil || r.lastAddress != invocation.Addr {
		r.crh = infra.NewCRH(invocation.Addr)
	}
	packet := common.NewRequestPacket(message)
	marshalledMessage, err := common.Marshall(*packet)
	if err != nil {return marshallingError(err)}
	if r.transportType == "tcp" {
		return r.SendReceiveTcp(marshalledMessage, message.IsReplyRequired())
	} else {
		return r.SendReceiveUdp(marshalledMessage, message.IsReplyRequired())
	}
}

func marshallingError(err error) *common.Termination {
	return &common.Termination{
		Result: fmt.Sprintf("%s: %s", "marshalling error: ", err),
	}
}

func (r *Requestor) SendReceiveTcp(message []byte, isReplyRequired bool) *common.Termination {
	r.crh.SendTcp(message)
	if isReplyRequired {
		data := r.crh.ReceiveTcp()
		return unpackTermination(data)
	} else {
		return &common.Termination{Result: "Success"}
	}
}

func unpackTermination(data []byte) *common.Termination {
	packet := &common.Packet{}
	err := common.Unmarshall(data, packet)
	if err != nil {return marshallingError(err)}
	ter := packet.Body.RepBody
	return &common.Termination{Result: ter}
}

func (r *Requestor) SendReceiveUdp(message []byte, isReplyRequired bool) *common.Termination {
	r.crh.SendUdp(message)
	if isReplyRequired {
		data := r.crh.ReceiveUdp()
		return unpackTermination(data)
	} else {
		return &common.Termination{Result: "Success"}
	}
}