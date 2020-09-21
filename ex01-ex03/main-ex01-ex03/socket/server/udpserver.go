package main_socket_server

import (
	"../../../protocol"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type UDPChatServer struct {
	listener *net.UDPConn
	clients []*udpclient
	mutex   *sync.Mutex
}

type udpclient struct {
	conn   *net.UDPConn
	addr   *net.UDPAddr
	name   string
}

func NewUDPServer() *UDPChatServer {
	return &UDPChatServer{
		mutex: &sync.Mutex{},
	}
}

func (s *UDPChatServer) deserializar(data []byte) (interface{}, error) {
	b := bytes.NewBuffer(data)
	// todo; what about other commands?
	cmd := protocol.MessageCommand{
		Name:    "",
		Message: "",
	}
	_, err := fmt.Fscanln(b, &cmd.Name, &cmd.Message)
	return cmd, err
}

func (s *UDPChatServer) serializar(cmd interface{}) ([]byte, error) {
	var b bytes.Buffer
	switch v := cmd.(type) {
	case protocol.MessageCommand:
		fmt.Fprintln(&b, v.Name, v.Message)
	default:
		log.Printf("Unknown client receiving command: %v", v)
	}

	return b.Bytes(), nil
}

func (s *UDPChatServer) serve(client *udpclient) {
	defer s.remove(client)
	for {
		buffer := make([]byte, 1024)
		_, addr, err := client.conn.ReadFromUDP(buffer)
		client.addr = addr
		if err != nil && err != io.EOF {
			log.Printf("Read error: %v", err)
		}
		cmd, err := s.deserializar(buffer)
		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.SendCommand:
				s.Broadcast(protocol.MessageCommand{
					Message: v.Message,
					Name:    client.name,
				})
			case protocol.NameCommand:
				client.name = v.Name
			case protocol.MessageCommand:
				s.Broadcast(v)
			}
		}
		if err == io.EOF {
			break
		}
	}
}

func (s *UDPChatServer) Start() {
	client := &udpclient{
		conn:   s.listener,
	}
	s.clients = append(s.clients, client)
	s.serve(client)
}

func (s *UDPChatServer) Listen(address string) error {
	addr,err := net.ResolveUDPAddr("udp",address)
	l, err := net.ListenUDP("udp", addr)
	if err == nil {
		s.listener = l
	}
	log.Printf("Listening on %v", address)
	return err
}

func (s *UDPChatServer) Close() {
	s.listener.Close()
}

func (s *UDPChatServer) Broadcast(command interface{}) error {
	for _, client := range s.clients {
		// TODO: handle error here?
		bytes, _ := s.serializar(command)
		client.conn.WriteToUDP(bytes, client.addr)
	}
	return nil
}


func (s *UDPChatServer) remove(client *udpclient) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// remove the connections from clients array
	for i, check := range s.clients {
		if check == client {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}
	log.Printf("Closing connection from %v", client.conn.RemoteAddr().String())
	client.conn.Close()
}


