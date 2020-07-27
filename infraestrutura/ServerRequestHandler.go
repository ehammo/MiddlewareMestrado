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
	mutex         *sync.Mutex
	uniqueId      int
}
type Client struct {
	conn        net.Conn
	udpconn     *net.UDPConn
	reader      *bufio.Reader
	Addr        *net.UDPAddr
	writer      *bufio.Writer
	UniqueId    int
	CurrentLane string
}

func NewServer(transportType string, address string) *ServerRequestHandler {
	var srh = &ServerRequestHandler {
		transportType: transportType,
		mutex: &sync.Mutex{},
		uniqueId: 0,
	}
	log.Printf("new server")
	srh.Listen(address)
	return srh
}

func (s *ServerRequestHandler) Listen(address string) {
	if s.transportType == "tcp" {
		l, err := net.Listen("tcp", address)
		failOnError(err, "error listening to address")
		s.listener = l
		log.Printf("Listening tcp on %v", address)
	} else if s.transportType == "udp" {
		addr,err := net.ResolveUDPAddr("udp",address)
		failOnError(err, "error resolving to address")
		l, err := net.ListenUDP("udp", addr)
		failOnError(err, "error listening to address")
		s.udplistener = l
		log.Printf("Listening udp on %v", address)
	}
}

func (s *ServerRequestHandler) AcceptNewClient() *Client {
	if s.transportType == "tcp" {
		conn, err := s.listener.Accept()
		failOnError(err, "Error accepting client")
		log.Printf("Accepting connection from %v", conn.RemoteAddr().String())
		s.mutex.Lock()
		var newClientId = s.uniqueId
		defer s.mutex.Unlock()
		client := &Client{
			conn:   conn,
			reader: bufio.NewReader(conn),
			writer: bufio.NewWriter(conn),
			UniqueId: newClientId,
			CurrentLane: "UNKNOWN",
		}
		log.Printf("id=%d",s.uniqueId)
		s.uniqueId = s.uniqueId+1
		return client
	} else {
		s.mutex.Lock()
		var newClientId = s.uniqueId
		defer s.mutex.Unlock()
		client := &Client {
			udpconn:   s.udplistener,
			reader: bufio.NewReader(s.udplistener),
			writer: bufio.NewWriter(s.udplistener),
			UniqueId: newClientId,
			CurrentLane: "UNKNOWN",
		}
		log.Println(s.uniqueId)
		s.uniqueId = s.uniqueId+1
		return client
	}
}

func (s *ServerRequestHandler) Receive(client *Client) ([]byte, *net.UDPAddr) {
	if s.transportType == "tcp" {
		cmd, err := client.reader.ReadBytes('\n')
		failOnError(err, "Read error:")
		return cmd, nil
	} else {
		log.Println("receiving udp")
		buffer := make([]byte, 1024)
		_, addr, err := client.udpconn.ReadFromUDP(buffer)
		failOnError(err, "Read error:")
		return buffer, addr
	}

}

func (s *ServerRequestHandler) Send(msg []byte, client *Client) {
	if s.transportType == "tcp" {
		_, err := client.writer.Write(msg)
		failOnError(err, "error writing")
		err = client.writer.Flush()
		failOnError(err, "error writing")
	} else {
		_, err := client.udpconn.WriteToUDP(msg, client.Addr)
		failOnError(err, "error writing")
	}
}
