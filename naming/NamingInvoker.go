package naming

import (
	c "../common"
	i "../infra"
	"bufio"
	"fmt"
)

type NamingInvoker struct {
	srh *i.ServerRequestHandler
	NamingImpl *NamingImpl
}

type Client struct {
	tcpReader *bufio.Reader
	tcpWriter *bufio.Writer
}

func NewNamingInvoker(address string) *NamingInvoker {
	return &NamingInvoker{
		srh: i.NewSRH(address),
		NamingImpl: NewNamingImpl(),
	}
}

func (n *NamingInvoker) Start() {
	fmt.Println("Starting naming invoker")
	for {
		conn := n.srh.AcceptNewClientTcp()
		newTcpClient := &Client{
			tcpReader: bufio.NewReader(*conn),
			tcpWriter: bufio.NewWriter(*conn),
		}
		go n.ServeTcp(newTcpClient)
	}
}

func (n *NamingInvoker) ServeTcp(client *Client) {
	for {
		data, err := n.srh.ReceiveTcp(client.tcpReader)
		if err != nil {
			fmt.Printf("Error receiving tcp data %s", err)
		}
		var packet = &c.Packet{}
		err = c.Unmarshall(data, packet)
		if err != nil {
			fmt.Printf("\nerro %s\n", err)
		}
		lookupMessage := c.CreateLookupMessageFromLookupPacket(packet)
		op := lookupMessage.Message.Operation
		topic := lookupMessage.Message.Topic
		fmt.Println(op)
		fmt.Println(topic)
		if op == "REGISTER" {
			n.NamingImpl.register(topic, lookupMessage.AOR)
		} else if op == "LOOKUP" {
			aor := n.NamingImpl.lookup(topic)
			packet := c.NewLookUpReplyPacket(aor)
			dataToSend, _ := c.Marshall(*packet)
			n.srh.SendTcp(dataToSend, client.tcpWriter)
		}

	}
}