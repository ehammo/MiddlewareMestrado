package infra

import (
	"bufio"
	"fmt"
	"net"
)

type ClientRequestHandler struct {
	reader    *bufio.Reader
	writer    *bufio.Writer
	udpconn   *net.UDPConn
	addr      *net.UDPAddr
	srvAddr   string
}

func NewCRH(address string) *ClientRequestHandler {
	return &ClientRequestHandler {
		srvAddr: address,
	}
}

func (c *ClientRequestHandler) DialTcp() {
	fmt.Println("Dialing")
	conn, err := net.Dial("tcp", c.srvAddr)
	failOnError(err, "error dialing address")
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)
}

func (c *ClientRequestHandler) DialUdp() {
	fmt.Println("Dialing")
	addr, err := net.ResolveUDPAddr("udp", c.srvAddr)
	failOnError(err, "error resolving address")
	conn, err := net.DialUDP("udp", nil, addr)
	failOnError(err, "error dialing address")
	c.udpconn = conn
}

func (c *ClientRequestHandler) ReceiveTcp() []byte {
	if c.reader == nil {
		c.DialTcp()
	}
	if c.reader != nil {
		buffer := make([]byte, 1024)
		size, err := c.reader.Read(buffer)
		cmd := buffer[:size]
		failOnError(err, "Error receiving message")
		if err == nil {
			return cmd
		}
	}
	return nil
}
func (c *ClientRequestHandler) ReceiveUdp() []byte {
	if c.udpconn == nil {
		c.DialUdp()
	}
	if c.udpconn != nil {
		buffer := make([]byte, 1024)
		_, addr, err := c.udpconn.ReadFromUDP(buffer)
		c.addr = addr
		failOnError(err, "Error receiving message")
		if err == nil {
			return buffer
		}
	}
	return nil
}

func (c *ClientRequestHandler) SendTcp(msg []byte) string {
	if c.writer == nil {
		c.DialTcp()
	}
	if c.writer != nil {
		_, err := c.writer.Write(msg)
		failOnError(err, "error writing")
		err = c.writer.Flush()
		failOnError(err, "error writing")
		if err != nil {
			return "error writing"
		}
	} else {
		fmt.Println("Error dialing")
		return "error dialing"
	}
	return "Send success"
}

func (c *ClientRequestHandler) SendUdp(msg []byte) string {
	if c.udpconn == nil {
		c.DialUdp()
	}
	if c.udpconn != nil {
		_, err := c.udpconn.Write(msg)
		failOnError(err, "error writing")
		if err != nil {
			return "error writing"
		}
	} else {
		return "error dialing"
	}
	return "Send success"
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
	}
}