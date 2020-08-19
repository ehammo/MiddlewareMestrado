package naming

import (
	c "../common"
	i "../infra"
	"bufio"
	"fmt"
	"reflect"
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
		op := packet.Body.ReqHeader.Operation
		body0 := packet.Body.ReqBody.Body[0]
		body1 := packet.Body.ReqBody.Body[1]
		service, _ := reflect.ValueOf(body0).Interface().(string)
		var aor c.AOR
		if body1 != nil {
			aor, _ = reflect.ValueOf(body1).Interface().(c.AOR)
		}
		fmt.Println("result:")
		fmt.Println(op)
		fmt.Println(service)
		if op == "Register" {
			n.NamingImpl.Register(service, &aor)
		} else if op == "LookUp" {
			aor := n.NamingImpl.Lookup(service)
			packet := c.NewReplyPacket(aor, "success")
			dataToSend, _ := c.Marshall(*packet)
			n.srh.SendTcp(dataToSend, client.tcpWriter)
		}

	}
}