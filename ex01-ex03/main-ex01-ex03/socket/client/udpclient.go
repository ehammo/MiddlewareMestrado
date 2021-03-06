package main_socket_client

import (
	"../../../protocol"
	"bytes"
	"fmt"
	"io"
	"time"
	"log"
	"net"
)

type UDPChatClient struct {
	conn      *net.UDPConn
	name      string
	incoming  chan protocol.MessageCommand
}

func NewUDPClient() *UDPChatClient {
	return &UDPChatClient{
		incoming: make(chan protocol.MessageCommand, 50000),
	}
}

func (c *UDPChatClient) Dial(address string) error {
	addr, err := net.ResolveUDPAddr("udp",address)
	conn, err := net.DialUDP("udp", nil, addr)

	if err == nil {
		c.conn = conn
	}

	return err
}


func (c *UDPChatClient) deserializar(data []byte) (interface{}, error) {
	b := bytes.NewBuffer(data)
	// todo; what about other commands?
	cmd := protocol.MessageCommand{
		Name:    "",
		Message: "",
	}
	_, err := fmt.Fscanln(b, &cmd.Name, &cmd.Message)
	return cmd, err
}

func (c *UDPChatClient) serializar(cmd interface{}) ([]byte, error) {
	var b bytes.Buffer
	switch v := cmd.(type) {
	case protocol.MessageCommand:
		fmt.Fprintln(&b, v.Name, v.Message)
	default:
		log.Printf("Unknown client receiving command: %v", v)
	}

	return b.Bytes(), nil
}

func (c *UDPChatClient) Start() {
	c.conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	for  {
		buffer := make([]byte, 1024)
		_, _, err := c.conn.ReadFromUDP(buffer)
		if err == io.EOF {
			println("EOF")
			// c.incoming <- io.EOF
			break
		} else if err != nil {
			if e, ok := err.(net.Error); !ok || !e.Timeout() {
				log.Printf("Read error %v", err)
			}
			break
		}
		cmd, err := c.deserializar(buffer)
		if cmd != nil {

			switch v := cmd.(type) {
			case protocol.MessageCommand:
				c.incoming <- v
			default:
				log.Printf("Unknown client receiving command: %v", v)
			}
		}
	}
	println("timeout. Closing channel")
	//close(c.incoming)
}

func (c *UDPChatClient) Close() {
	c.conn.Close()
}

func (c *UDPChatClient) Incoming() chan protocol.MessageCommand {
	return c.incoming
}

func (c *UDPChatClient) Send(command interface{}) error {
	cmdBytes, _ := c.serializar(command)
	_, err := c.conn.Write(cmdBytes)
	return err
}

func (c *UDPChatClient) SendMessage(message string) error {
	return c.Send(protocol.MessageCommand{
		Name: c.name,
		Message: message,
	})
}

func (c *UDPChatClient) SetName(name string) error {
	c.name = name
	return nil
}

func (c *UDPChatClient) Clean() {
	close(c.incoming)
}