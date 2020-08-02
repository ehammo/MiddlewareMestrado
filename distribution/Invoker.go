package distribution

import (
	common "../common"
	infra "../infra"
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Invoker struct {
	srh           *infra.ServerRequestHandler
	transportType string
	mutex         sync.Mutex
	clients       map[int]*Client
	uniqueId      int
}

type Client struct {
	currentLane string
	tcpReader   *bufio.Reader
	tcpWriter   *bufio.Writer
	udpAddr     *net.UDPAddr
	id          int
}

func NewInvoker(address string, transportType string) *Invoker {
	return &Invoker{
		srh: infra.NewSRH(address),
		transportType: transportType,
		clients: make(map[int]*Client),
		uniqueId: 0,
	}
}

func (i *Invoker) Start() {
	if i.transportType == "tcp" {
		for {
			conn := i.srh.AcceptNewClientTcp()
			newTcpClient := &Client{
				tcpReader: bufio.NewReader(*conn),
				tcpWriter: bufio.NewWriter(*conn),
				id:        i.uniqueId,
			}
			i.addClientOnList(newTcpClient)
			go i.ServeTcp(newTcpClient)
		}
	} else {
		i.ServeUdp()
	}
}

func (i *Invoker) ServeUdp() {
	for {
		data, addr, err := i.srh.ReceiveUdp()
		if err != nil {
			fmt.Printf("Error receiving udp data %s", err)
		}
		newUdpClient := i.findAddUdpClient(addr)
		i.unmarshallAndRun(data, newUdpClient)
	}
}

func (i *Invoker) ServeTcp(client *Client) {
	fmt.Println("Serving client C"+strconv.Itoa(client.id))
	for {
		data, err := i.srh.ReceiveTcp(client.tcpReader)
		if err != nil {
			fmt.Printf("Error receiving tcp data %s", err)
		}
		i.unmarshallAndRun(data, client)
	}
}

func (i *Invoker) unmarshallAndRun(data []byte, client *Client){
	fmt.Println("unmarshalling and running message from C"+strconv.Itoa(client.id))
	message := &common.Message{}
	var err = common.Unmarshall(data, message)
	if err != nil {
		fmt.Printf("Error unmarshelling %s", err)
	}
	fmt.Println("message = "+message.Operation+", "+message.Topic)
	i.runCmd(client, message)
}

func (i *Invoker) runCmd(c *Client, message *common.Message) {
	fmt.Println("running command "+message.Operation+" from client C"+strconv.Itoa(c.id))
	if message.Operation == "REGISTER" || message.Operation == "LANE" {
		c.currentLane = message.Topic
		i.mutex.Lock()
		i.clients[c.id] = c
		i.mutex.Unlock()
	} else if message.Operation == "BREAK" {
		var lane = message.Topic
		i.mutex.Lock()
		for _, client := range i.clients {
			if client != nil && strings.Contains(lane, client.currentLane) {
				if i.transportType == "tcp" {
					data, err := common.Marshall(*message)
					if err != nil {
						fmt.Printf("Error marshelling %s", err)
					}
					err = i.srh.SendTcp(data, client.tcpWriter)
					if err != nil {
						fmt.Printf("Error sending %s", err)
					}
				} else {
					data, err := common.Marshall(*message)
					if err != nil {
						fmt.Printf("Error marshelling %s", err)
					}
					err = i.srh.SendUdp(data, client.udpAddr)
					if err != nil {
						fmt.Printf("Error sending %s", err)
					}
				}
			}
		}
		i.mutex.Unlock()
	}
}

func (i *Invoker) addClientOnList(newClient *Client) {
	i.mutex.Lock()
	i.clients[i.uniqueId] = newClient
	i.uniqueId = i.uniqueId + 1
	i.mutex.Unlock()
}

func (i *Invoker) findAddUdpClient(addr *net.UDPAddr) *Client {
	newUdpClient := &Client{
		udpAddr:   addr,
		id:        i.uniqueId,
	}
	var found = false
	for _, client := range i.clients {
		if client.udpAddr == addr {
			found = true
			newUdpClient = client
		}
	}
	if !found {
		i.addClientOnList(newUdpClient)
	}
	return newUdpClient
}