package naming

import (
	"crypto/rsa"
	"fmt"

	c "../common"
	d "../distribution"
	i "../infra"
)

type NamingProxy struct {
	address string
	crh     *i.ClientRequestHandler
}

func NewNamingProxy(address string) *NamingProxy {
	return &NamingProxy{
		address: address,
		crh:     i.NewCRH(address),
	}
}

func (n *NamingProxy) Register(service string, aor *c.AOR) {
	message := &c.Message{
		Operation: "Register",
		Topic:     service,
		AOR:       aor,
	}
	packet := c.NewRequestPacket(message)
	data, _ := c.Marshall(*packet)
	n.crh.SendTcp(data)
}

func (n *NamingProxy) RegisterKey(service string, key *rsa.PublicKey) {
	message := &c.Message{
		Operation: "RegisterKey",
		Topic:     service,
		AOR:       key,
	}
	packet := c.NewRequestPacket(message)
	data, _ := c.Marshall(*packet)
	n.crh.SendTcp(data)
}

func (n *NamingProxy) LookUp(service string) *c.AOR {
	fmt.Println("Looking up")
	message := &c.Message{
		Operation: "LookUp",
		Topic:     service,
	}
	invocation := &c.Invocation{
		Addr:    n.address,
		Message: message,
	}
	requestor := d.NewRequestor("tcp")
	ter := requestor.Invoke(invocation)
	fmt.Println(ter)
	aor, _ := ter.Result.(*c.AOR)
	return aor
}

func (n *NamingProxy) LookUpKey(service string) *rsa.PublicKey {
	fmt.Println("Looking up")
	message := &c.Message{
		Operation: "LookUpKey",
		Topic:     service,
	}
	invocation := &c.Invocation{
		Addr:    n.address,
		Message: message,
	}
	requestor := d.NewRequestor("tcp")
	ter := requestor.Invoke(invocation)
	key, _ := ter.Result.(*rsa.PublicKey)
	return key
}
