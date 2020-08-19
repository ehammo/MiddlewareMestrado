package naming

import (
	c "../common"
	d "../distribution"
	i "../infra"
	"fmt"
	"reflect"
)

type NamingProxy struct {
	address string
	crh     *i.ClientRequestHandler
}

func NewNamingProxy() *NamingProxy {
	address := "localhost:1243"
	return &NamingProxy{
		address: address,
		crh:      i.NewCRH(address),
	}
}

func (n *NamingProxy) Register(service string, aor *c.AOR) {
	message := c.Message{
		Operation: "REGISTER",
		Topic:     service,
		AOR:       aor,
	}
	packet := *c.NewRequestPacket(message)
	data, _ := c.Marshall(packet)
	n.crh.SendTcp(data)
}

func (n *NamingProxy) LookUp(service string) *c.AOR {
	fmt.Println("Looking up")
	message := &c.Message{
		Operation: "LOOKUP",
		Topic:     service,
	}
	invocation := &c.Invocation{
		Addr:    n.address,
		Message: message,
	}
	requestor := d.NewRequestor("tcp")
	ter := requestor.Invoke(invocation)
	aor,_ := reflect.ValueOf(ter.Result).Interface().(c.AOR)
	return &aor
}