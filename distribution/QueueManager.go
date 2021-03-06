package distribution

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"sync"

	common "../common"
	infra "../infra"
)

type QueueManager struct {
	srh           *infra.ServerRequestHandler
	transportType string
	mutex         *sync.Mutex
	clients       map[int]*Client
	uniqueId      int
	kp            *common.Keypair
}

type Client struct {
	currentLane string
	EventBus    *EventBus
	tcpReader   *bufio.Reader
	tcpWriter   *bufio.Writer
	udpAddr     *net.UDPAddr
	Id          int
	i           int
	Pub         *rsa.PublicKey
}

func NewQueueManager(address string, transportType string, kp *common.Keypair) *QueueManager {
	return &QueueManager{
		srh:           infra.NewSRH(address),
		mutex:         &sync.Mutex{},
		transportType: transportType,
		clients:       make(map[int]*Client),
		uniqueId:      0,
		kp:            kp,
	}
}

func (qm *QueueManager) Start() {
	if qm.transportType == "tcp" {
		for {
			conn := qm.srh.AcceptNewClientTcp()
			var eventBus = NewEventBus()
			newTcpClient := &Client{
				tcpReader: bufio.NewReader(*conn),
				tcpWriter: bufio.NewWriter(*conn),
				Id:        -1,
				i:         qm.uniqueId,
				EventBus:  eventBus,
			}
			eventBus.SetQueueManager(qm)
			eventBus.SetClient(newTcpClient)
			qm.addClientOnList(newTcpClient)
			go qm.ServeTcp(newTcpClient)
		}
	} else {
		qm.ServeUDP()
	}
}

func (qm *QueueManager) ServeUDP() {
	for {
		data, addr, err := qm.srh.ReceiveUDP()
		if err != nil {
			fmt.Printf("Error receiving udp data %s", err)
		}
		newUDPClient := qm.findAddUDPClient(addr)
		qm.unmarshallAndRun(data, newUDPClient)
	}
}

func (qm *QueueManager) ServeTcp(client *Client) {
	fmt.Println("Serving client C" + strconv.Itoa(client.i))
	errTotal := 0
	for {
		data, err := qm.srh.ReceiveTcp(client.tcpReader)
		if err != nil {
			errTotal++
			fmt.Printf("Error receiving tcp data %s", err)
		} else {
			qm.unmarshallAndRun(data, client)
		}
		if errTotal > 5 {
			break
		}
	}
}

func (qm *QueueManager) unmarshallAndRun(data []byte, client *Client) {
	uncriptedData := common.Decrypt(data, qm.kp.Priv, false)
	fmt.Println("unmarshalling and running message from C" + strconv.Itoa(client.i))
	packet := &common.Packet{}
	var err = common.Unmarshall(uncriptedData, packet)
	if err != nil {
		fmt.Printf("Error unmarshelling %s", err)
	}
	fmt.Println("Unmarshalled")
	qm.runCmd(client, packet)
}

func (qm *QueueManager) runCmd(c *Client, packet *common.Packet) {
	clientID := packet.Header.ClientId
	fmt.Println("Message from C", clientID)
	c.Id = clientID
	operation := packet.Body.ReqHeader.Operation
	fmt.Println("op", operation)
	fmt.Println("checking body length")
	if len(packet.Body.ReqBody.Body) > 0 {
		body := packet.Body.ReqBody.Body[0]
		message := &common.Message{
			Operation: operation,
			Topic:     body,
		}
		if c.Pub == nil {
			c.Pub = qm.findPubWithClientId(c.Id)
		}
		fmt.Println("running command " + message.Operation + " from client C" + strconv.Itoa(c.Id))
		if message.Operation == "Register" {
			c.EventBus.RegisterOnLane(message.Topic.(string))
		} else if message.Operation == "RegisterKey" {
			pubMap := message.Topic.(map[string]interface{})
			Nstring := pubMap["N"].(string)
			Eint := int(pubMap["E"].(float64))
			fmt.Println(strconv.Itoa(Eint))
			fmt.Println(Nstring)
			N := big.Int{}
			N.SetString(Nstring, 10)
			clientPub := &rsa.PublicKey{
				N: &N,
				E: Eint,
			}
			c.Pub = clientPub
		} else if message.Operation == "ChangeLane" {
			c.EventBus.ChangeLane(message.Topic.(string))
		} else if message.Operation == "BroadcastEvent" {
			c.EventBus.BroadcastEvent(message.Topic.(string))
		}
		fmt.Println("Finishing running command " + message.Operation + " from client C" + strconv.Itoa(c.Id))
	}
}

func (qm *QueueManager) sendMessage(message *common.Message, client *Client) {
	fmt.Println("Sending message")
	packet := common.NewRequestPacket(message)
	data, err := common.Marshall(*packet)
	if err != nil {
		fmt.Printf("Error marshelling %s", err)
	}
	criptedData := common.Encrypt(data, client.Pub, true)
	if qm.transportType == "tcp" {
		err = qm.srh.SendTcp(criptedData, client.tcpWriter)
		if err != nil {
			fmt.Printf("Error sending %s", err)
		}
	} else {
		err = qm.srh.SendUDP(data, client.udpAddr)
		if err != nil {
			fmt.Printf("Error sending %s", err)
		}
	}
}

func (qm *QueueManager) addClientOnList(newClient *Client) {
	qm.mutex.Lock()
	fmt.Println("add client on list. Locked")
	qm.clients[qm.uniqueId] = newClient
	qm.uniqueId = qm.uniqueId + 1
	qm.mutex.Unlock()
	fmt.Println("add client on list. Unlocked")
}

func (qm *QueueManager) findPubWithClientId(id int) *rsa.PublicKey {
	qm.mutex.Lock()
	fmt.Println("findPubWithClientId. Locked")
	var pub *rsa.PublicKey
	for _, client := range qm.clients {
		if client.Id == id {
			pub = client.Pub
		}
	}
	qm.mutex.Unlock()
	fmt.Println("findPubWithClientId. Unlocked")
	return pub
}

func (qm *QueueManager) findAddUDPClient(addr *net.UDPAddr) *Client {
	eventBus := NewEventBus()
	newUDPClient := &Client{
		udpAddr:  addr,
		Id:       qm.uniqueId,
		EventBus: eventBus,
	}
	eventBus.SetQueueManager(qm)
	eventBus.SetClient(newUDPClient)
	var found = false
	for _, client := range qm.clients {
		if client.udpAddr == addr {
			found = true
			newUDPClient = client
		}
	}
	if !found {
		qm.addClientOnList(newUDPClient)
	}
	return newUDPClient
}
