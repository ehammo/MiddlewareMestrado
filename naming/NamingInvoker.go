package naming

import (
	"bufio"
	"crypto/rsa"
	"fmt"

	c "../common"
	i "../infra"
)

type NamingInvoker struct {
	srh        *i.ServerRequestHandler
	NamingImpl *NamingImpl
}

type Client struct {
	tcpReader *bufio.Reader
	tcpWriter *bufio.Writer
}

func NewNamingInvoker(address string) *NamingInvoker {
	return &NamingInvoker{
		srh:        i.NewSRH(address),
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
	var totalErr = 0
	for {
		data, err := n.srh.ReceiveTcp(client.tcpReader)
		fmt.Println("Naming server received something")
		if err != nil {
			fmt.Printf("Error receiving tcp data %s", err)
			totalErr++
			if totalErr == 5 {
				break
			}
		} else {
			var packet = &c.Packet{}
			err = c.Unmarshall(data, packet)
			if err != nil {
				totalErr++
				if totalErr == 5 {
					break
				}
				fmt.Printf("\nMarshalling error %s\n", err)
			} else {
				op := packet.Body.ReqHeader.Operation
				body0 := packet.Body.ReqBody.Body[0]
				body1 := packet.Body.ReqBody.Body[1]
				service, _ := body0.(string)
				fmt.Println("result:", op, service)
				var aor *c.AOR
				var key *rsa.PublicKey
				if body1 != nil && op == "Register" {
					fmt.Println("body1 is not null")
					add := body1.(map[string]interface{})["Address"].(string)
					fmt.Println(add)
					objId := body1.(map[string]interface{})["ObjectId"].(string)
					fmt.Println(objId)
					E := interface{}(body1.(map[string]interface{})["E"]).(float64)
					p := body1.(map[string]interface{})["Protocol"].(string)
					fmt.Println(p)
					N := body1.(map[string]interface{})["N"].(string)
					fmt.Println(N)
					aor = &c.AOR{
						Address:  add,
						E:        int(E),
						N:        N,
						ObjectId: objId,
						Protocol: p,
					}
					fmt.Println("aor: ", aor.ToString())
				}
				if op == "Register" {
					n.NamingImpl.Register(service, aor)
				} else if op == "RegisterKey" {
					n.NamingImpl.RegisterKey(service, key)
				} else if op == "LookUp" {
					aor := n.NamingImpl.Lookup(service)
					fmt.Println("Creating packet")
					packet := c.NewReplyPacket(aor, "success")
					fmt.Println("marshalling")
					dataToSend, _ := c.Marshall(*packet)
					fmt.Println("prep to send")
					n.srh.SendTcp(dataToSend, client.tcpWriter)
					fmt.Println("sent")
				}
			}
		}
	}
}
