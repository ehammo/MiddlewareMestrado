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
	marshalledMessage, err := common.Marshall(*invocation.Message)
	if err != nil {
		return &common.Termination{
			Result: fmt.Sprintf("%s: %s", "marshalling error: ", err),
		}
	}
	if r.transportType == "tcp" {
		return &common.Termination{
			Result: r.crh.SendTcp(marshalledMessage),
		}
	} else {
		return &common.Termination{
			Result: r.crh.SendUdp(marshalledMessage),
		}
	}
}