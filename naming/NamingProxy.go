package naming

import (
	c "../common"
	i "../infra"
	"encoding/json"
	"fmt"
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
	}
	packet := *c.NewLookUpRequestPacket(message, aor)
	data, _ := c.Marshall(packet)
	n.crh.SendTcp(data)
}

func (n *NamingProxy) LookUp(service string) *c.AOR {
	fmt.Println("Looking up")
	message := c.Message{
		Operation: "LOOKUP",
		Topic:     service,
	}
	packet := *c.NewLookUpRequestPacket(message, nil)
	data, _ := c.Marshall(packet)
	fmt.Println("sending data")
	n.crh.SendTcp(data)
	received := n.crh.ReceiveTcp()
	var replyPacket = &c.Packet{}
	c.Unmarshall(received, replyPacket)
	var aor = &c.AOR{}
	json.Unmarshal(replyPacket.Body, aor)
	return aor
}