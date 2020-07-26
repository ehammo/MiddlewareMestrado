package infraestrutura

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type ClientRequestHandler struct {
	transportType string
	reader *bufio.Reader
	writer *bufio.Writer
	udpconn *net.UDPConn
	addr      *net.UDPAddr
}

func NewClient(transportType string) *ClientRequestHandler {
	return &ClientRequestHandler {
		transportType: transportType,
	}
}

func (c *ClientRequestHandler) Dial(address string) {
	if c.transportType == "tcp" {
		conn, err := net.Dial("tcp", address)
		failOnError(err, "error dialing address")
		c.reader = bufio.NewReader(conn)
		c.writer = bufio.NewWriter(conn)
	} else if c.transportType == "udp" {
		addr, err := net.ResolveUDPAddr("udp",address)
		failOnError(err, "error resolving address")
		conn, err := net.DialUDP("udp", nil, addr)
		failOnError(err, "error dialing address")
		c.udpconn = conn
	} else {
		failOnError(nil, "invalid transport type")
	}
}

func (c *ClientRequestHandler) Receive() []byte {
	if c.transportType == "tcp" {
		cmd, err := c.reader.ReadBytes('\n')
		failOnError(err, "Error receiving message")
		if err == nil {
			return cmd
		} else {
			return nil
		}
	} else {
		buffer := make([]byte, 1024)
		_, addr, err := c.udpconn.ReadFromUDP(buffer)
		c.addr = addr
		if err == nil {
			return buffer
		} else {
			return nil
		}
	}

}

func (c *ClientRequestHandler) Send(msg []byte) {
	if c.transportType == "tcp" {
		_, err := c.writer.Write(msg)
		failOnError(err, "error writing")
		err = c.writer.Flush()
		failOnError(err, "error writing")
	} else {
		log.Println(string(msg))
		_, err := c.udpconn.Write(msg)
		failOnError(err, "error writing")
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
	}
}