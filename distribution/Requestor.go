package distribution

import (
	"fmt"

	common "../common"
	infra "../infra"
)

type Requestor struct {
	crh           *infra.ClientRequestHandler
	transportType string
	lastAddress   string
}

func NewRequestor(transportType string) *Requestor {
	return &Requestor{
		transportType: transportType,
	}
}

func (r *Requestor) Invoke(invocation *common.Invocation) *common.Termination {
	message := invocation.Message
	fmt.Println("Requestor: invoking " + message.Operation)
	if r.crh == nil || r.lastAddress != invocation.Addr {
		r.crh = infra.NewCRH(invocation.Addr)
	}
	packet := common.NewRequestPacket(message)
	marshalledMessage, err := common.Marshall(*packet)
	if err != nil {
		return marshallingError(err)
	}
	if r.transportType == "tcp" {
		return r.SendReceiveTcp(marshalledMessage, message.IsReplyRequired())
	} else {
		return r.SendReceiveUDP(marshalledMessage, message.IsReplyRequired())
	}
}

func marshallingError(err error) *common.Termination {
	return &common.Termination{
		Result: fmt.Sprintf("%s: %s", "Marshalling error: ", err),
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
	if err != nil {
		return marshallingError(err)
	}
	body1 := packet.Body.RepBody.Body[0]
	add := body1.(map[string]interface{})["Address"].(string)
	objId := body1.(map[string]interface{})["ObjectId"].(string)
	E := interface{}(body1.(map[string]interface{})["E"]).(float64)
	p := body1.(map[string]interface{})["Protocol"].(string)
	N := body1.(map[string]interface{})["N"].(string)
	aor := &common.AOR{
		Address:  add,
		E:        int(E),
		N:        N,
		ObjectId: objId,
		Protocol: p,
	}
	return &common.Termination{Result: aor}
}

func (r *Requestor) SendReceiveUDP(message []byte, isReplyRequired bool) *common.Termination {
	r.crh.SendUDP(message)
	if isReplyRequired {
		data := r.crh.ReceiveUDP()
		return unpackTermination(data)
	} else {
		return &common.Termination{Result: "Success"}
	}
}
