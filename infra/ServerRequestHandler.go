package infra

import (
	"bufio"
	"log"
	"net"
)

type ServerRequestHandler struct {
	listener      net.Listener
	udplistener   *net.UDPConn
	address       string
}


func NewSRH(address string) *ServerRequestHandler {
	return &ServerRequestHandler {
		address: address,
	}
}

func (s *ServerRequestHandler) ListenTcp() {
	l, err := net.Listen("tcp", s.address)
	failOnError(err, "error listening to address")
	s.listener = l
	log.Printf("Listening tcp on %v", s.address)
}

func (s *ServerRequestHandler) ListenUdp() {
	addr,err := net.ResolveUDPAddr("udp",s.address)
	failOnError(err, "error resolving to address")
	l, err := net.ListenUDP("udp", addr)
	failOnError(err, "error listening to address")
	s.udplistener = l
	log.Printf("Listening udp on %v", s.address)
}

func (s *ServerRequestHandler) AcceptNewClientTcp() *net.Conn {
	if s.listener == nil {
		s.ListenTcp()
	}
	conn, err := s.listener.Accept()
	failOnError(err, "Error accepting client")
	log.Printf("Accepting connection from %v", conn.RemoteAddr().String())
	return &conn
}

func (s *ServerRequestHandler) getUdpConn() *net.UDPConn {
	if s.udplistener == nil {
		s.ListenUdp()
	}
	return s.udplistener
}

func (s *ServerRequestHandler) ReceiveTcp(reader *bufio.Reader) ([]byte, error) {
	buffer := make([]byte, 1024)
	size, err := reader.Read(buffer)
	cmd := buffer[:size]
	failOnError(err, "Read error:")
	return cmd, err
}
func (s *ServerRequestHandler) ReceiveUdp() ([]byte, *net.UDPAddr, error) {
	buffer := make([]byte, 1024)
	size, addr, err := s.udplistener.ReadFromUDP(buffer)
	failOnError(err, "Read error:")
	cmd := buffer[:size]
	return cmd, addr, err
}

func (s *ServerRequestHandler) SendTcp(msg []byte, writer *bufio.Writer) error {
	_, err := writer.Write(msg)
	failOnError(err, "error writing")
	err = writer.Flush()
	failOnError(err, "error writing")
	return err
}

func (s *ServerRequestHandler) SendUdp(msg []byte, addr *net.UDPAddr) error {
	_, err := s.udplistener.WriteToUDP(msg, addr)
	failOnError(err, "error writing")
	return err
}