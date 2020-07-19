package main_rpc

import (
	"../../protocol"
	pb "../../protocol/rpc"
	clientI "../socket/client"
	"context"
	"errors"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

type RpcClient struct {
	client pb.RpcChatClient
	conn *grpc.ClientConn
	stream pb.RpcChat_SendMessagesClient
	name      string
	incoming  chan protocol.MessageCommand
}

func NewRpcClient() clientI.ChatClient {
	return &RpcClient{
		incoming: make(chan protocol.MessageCommand, 50000),
	}
}

func (c *RpcClient) SetName(name string) error{
	println("Setting my name as ", name)
	c.name = name
	return nil
}

func (c *RpcClient) Dial(address string) error {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	println("Dialing ", address)
	conn, err := grpc.Dial(address, opts...)
	if err == nil {
		println("Dialed successfully")
		c.conn = conn
		c.client = pb.NewRpcChatClient(conn)
	} else {
		log.Printf("dial error=%s",err.Error())
	}

	return err
}

func (c *RpcClient) Send(command interface{}) error {
	var in pb.Messages
	switch v := command.(type) {
	case protocol.MessageCommand:
		in = pb.Messages{
			Name:                 v.Name,
			Message:              v.Message,
		}
	}
	return c.stream.Send(&in)
}

func (c *RpcClient) SendMessage(message string) error {
	if c.stream != nil {
		serializedMessage := c.serialize(c.name, message)
		err := c.stream.Send(&serializedMessage)
		return err
	}
	return errors.New("no stream")
}


func (c *RpcClient) serialize(name string, message string) pb.Messages {
	return pb.Messages{
		Name:                 name,
		Message:              message,
	}
}

func (c *RpcClient) messageToCommand(in *pb.Messages) protocol.MessageCommand {
	return protocol.MessageCommand{
		Name:    in.Name,
		Message: in.Message,
	}
}

func (c *RpcClient) Close() {
	c.conn.Close()
}

func (c *RpcClient) Start() {
	println("Sending stuff")
	ctx, cancel := context.WithTimeout(context.Background(), 30000*time.Second)
	c.stream, _ = c.client.SendMessages(ctx)
	for {
		in, err := c.stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Read error %v", err)
			break
		}
		if in != nil {
			c.incoming <- c.messageToCommand(in)
		}
	}
	defer cancel()
}

func (c *RpcClient) Incoming() chan protocol.MessageCommand {
	return c.incoming
}

func (c *RpcClient) Clean() {
	c.incoming = make(chan protocol.MessageCommand)
}