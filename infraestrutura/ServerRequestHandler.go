package infraestrutura

import (
	"bufio"
	"log"
	"net"
	"sync"
)

type ServerRequestHandler struct {
	transportType string
	listener      net.Listener
	udplistener   *net.UDPConn
	clients       []*client
	mutex         *sync.Mutex
	uniqueId      int
	maxClient     int
}
type client struct {
	conn      net.Conn
	udpconn   *net.UDPConn
	reader    *bufio.Reader
	addr      *net.UDPAddr
	name      string
	writer    *bufio.Writer
	uniqueId  int
}

func NewServer(transportType string) *ServerRequestHandler {
	return &ServerRequestHandler {
		transportType: transportType,
		mutex: &sync.Mutex{},
		uniqueId: 0,
		maxClient: 10,
		clients: make([]*client, 10),
	}
}

func (s *ServerRequestHandler) Listen(address string) {
	if s.transportType == "tcp" {
		l, err := net.Listen("tcp", address)
		failOnError(err, "error listening to address")
		s.listener = l
		log.Printf("Listening on %v", address)
	} else if s.transportType == "udp" {
		addr,err := net.ResolveUDPAddr("udp",address)
		failOnError(err, "error resolving to address")
		l, err := net.ListenUDP("udp", addr)
		failOnError(err, "error listening to address")
		s.udplistener = l
		log.Printf("Listening on %v", address)
	}
}

func (s *ServerRequestHandler) AcceptNewClient() int {
	if s.transportType == "tcp" {
		conn, err := s.listener.Accept()
		failOnError(err, "Error accepting client")
		log.Printf("Accepting connection from %v, total clients: %v", conn.RemoteAddr().String(), len(s.clients)+1)
		s.mutex.Lock()
		var newClientId = s.uniqueId
		defer s.mutex.Unlock()
		client := &client{
			conn:   conn,
			reader: bufio.NewReader(conn),
			writer: bufio.NewWriter(conn),
			uniqueId: newClientId,
		}
		log.Printf("id=%d",s.uniqueId)
		s.clients[s.uniqueId] = client
		s.uniqueId = s.uniqueId+1
		return newClientId
	} else {
		if s.uniqueId >= 1 {
			return -1
		}
		s.mutex.Lock()
		var newClientId = s.uniqueId
		defer s.mutex.Unlock()
		client := &client {
			udpconn:   s.udplistener,
			reader: bufio.NewReader(s.udplistener),
			writer: bufio.NewWriter(s.udplistener),
			uniqueId: newClientId,
		}
		log.Println(s.uniqueId)
		s.clients[s.uniqueId] = client
		s.uniqueId = s.uniqueId+1
		return newClientId
	}
}

func (s *ServerRequestHandler) Receive(clientId int) []byte {
	if s.transportType == "tcp" {
		cmd, err := s.clients[clientId].reader.ReadBytes('\n')
		failOnError(err, "Read error:")
		return cmd
	} else {
		log.Println("receiving udp")
		buffer := make([]byte, 1024)
		_, addr, err := s.clients[clientId].udpconn.ReadFromUDP(buffer)
		s.clients[clientId].addr = addr
		failOnError(err, "Read error:")
		return buffer
	}

}

func (s *ServerRequestHandler) Send(msg []byte, clientId int) {
	if s.transportType == "tcp" {
		_, err := s.clients[clientId].writer.Write(msg)
		failOnError(err, "error writing")
		err = s.clients[clientId].writer.Flush()
		failOnError(err, "error writing")
	} else {
		_, err := s.clients[clientId].udpconn.WriteToUDP(msg, s.clients[clientId].addr)
		failOnError(err, "error writing")
	}
}
