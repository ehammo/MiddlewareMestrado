package distribution

import (
	common "../common"
	infra "../infra"
	"bufio"
	"fmt"
	"net"
	"strconv"
	"sync"
)

type Invoker struct {
	srh           *infra.ServerRequestHandler
	transportType string
	mutex         *sync.Mutex
	clients       map[int]*Client
	uniqueId      int
}

type Client struct {
	currentLane string
	EventBus    *EventBus
	tcpReader   *bufio.Reader
	tcpWriter   *bufio.Writer
	udpAddr     *net.UDPAddr
	id          int
}

func NewInvoker(address string, transportType string) *Invoker {
	return &Invoker{
		srh: infra.NewSRH(address),
		mutex: &sync.Mutex{},
		transportType: transportType,
		clients: make(map[int]*Client),
		uniqueId: 0,
	}
}

func (i *Invoker) Start() {
	if i.transportType == "tcp" {
		for {
			conn := i.srh.AcceptNewClientTcp()
			var eventBus = NewEventBus()
			newTcpClient := &Client{
				tcpReader: bufio.NewReader(*conn),
				tcpWriter: bufio.NewWriter(*conn),
				id:        i.uniqueId,
				EventBus:  eventBus,
			}
			eventBus.SetInvoker(i)
			eventBus.SetClient(newTcpClient)
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
		reply := common.NewReplyPacket("Success")
		dataToSend, err := common.Marshall(*reply)
		//marshallerror
		i.srh.SendUdp(dataToSend, newUdpClient.udpAddr)
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
		reply := common.NewReplyPacket("Success")
		dataToSend, err := common.Marshall(*reply)
		//marshallerror
		i.srh.SendTcp(dataToSend, client.tcpWriter)
	}
}

func (i *Invoker) unmarshallAndRun(data []byte, client *Client){
	fmt.Println("unmarshalling and running message from C"+strconv.Itoa(client.id))
	packet := &common.Packet{}
	var err = common.Unmarshall(data, packet)
	if err != nil {
		fmt.Printf("Error unmarshelling %s", err)
	}
	i.runCmd(client, packet)
}

func (i *Invoker) runCmd(c *Client, packet *common.Packet) {
	message := &common.Message{
		Operation: string(packet.Header),
		Topic:     string(packet.Body),
	}
	fmt.Println("running command "+message.Operation+" from client C"+strconv.Itoa(c.id))
	if message.Operation == "REGISTER" {
		c.EventBus.RegisterOnLane(message.Topic)
	} else if message.Operation == "LANE" {
		c.EventBus.ChangeLane(message.Topic)
	} else if message.Operation == "BREAK" {
		c.EventBus.BroadcastEvent(message.Topic)
	}
}

func (i *Invoker) sendMessage(message *common.Message, client *Client) {
	fmt.Println("Sending message")
	packet := common.NewRequestPacket(*message)
	if i.transportType == "tcp" {
		data, err := common.Marshall(*packet)
		if err != nil {
			fmt.Printf("Error marshelling %s", err)
		}
		err = i.srh.SendTcp(data, client.tcpWriter)
		if err != nil {
			fmt.Printf("Error sending %s", err)
		}
	} else {
		data, err := common.Marshall(*packet)
		if err != nil {
			fmt.Printf("Error marshelling %s", err)
		}
		err = i.srh.SendUdp(data, client.udpAddr)
		if err != nil {
			fmt.Printf("Error sending %s", err)
		}
	}
}

func (i *Invoker) addClientOnList(newClient *Client) {
	i.mutex.Lock()
	i.clients[i.uniqueId] = newClient
	i.uniqueId = i.uniqueId + 1
	i.mutex.Unlock()
}

func (i *Invoker) findAddUdpClient(addr *net.UDPAddr) *Client {
	eventBus := NewEventBus()
	newUdpClient := &Client{
		udpAddr:   addr,
		id:        i.uniqueId,
		EventBus:  eventBus,
	}
	eventBus.SetInvoker(i)
	eventBus.SetClient(newUdpClient)
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