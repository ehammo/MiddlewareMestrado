package main_rpc

import (
	protocol "../../protocol"
	pb "../../protocol/rpc"
	server "../socket/server"
	"flag"
	grpc "google.golang.org/grpc"
	"io"
	"log"
	"net"
	"sync"
)

type RpcServer struct {
	unimplementedServer *pb.RpcChatServer
	listener net.Listener
	clients map[string]*client
	mutex   *sync.Mutex
}
type client struct {
	stream pb.RpcChat_SendMessagesServer
}

func NewRpcServer() server.ChatServer {
	return &RpcServer{
		mutex: &sync.Mutex{},
		clients: make(map[string]*client),
	}
}

func (s *RpcServer) SendMessages(stream pb.RpcChat_SendMessagesServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			println("eof")
			return nil
		}
		if err != nil {
			println("erro no recebimento "+err.Error())
			return err
		}

		client := client{
			stream: stream,
		}
		s.mutex.Lock()
		s.clients[in.Name] = &client
		println("server received:",in.Name, in.Message)
		cmd := protocol.MessageCommand{
			Name:    in.Name,
			Message: in.Message,
		}
		err = s.Broadcast(cmd)
		if err != nil {
			println("erro no broadcast "+err.Error())
			return err
		}
		s.mutex.Unlock()
	}
}

func (s *RpcServer) Start() {
	grpcServer := grpc.NewServer()
	pb.RegisterRpcChatServer(grpcServer, s)
	grpcServer.Serve(s.listener)
}

func (s *RpcServer) Listen(address string) error {
	flag.Parse()
	lis, err := net.Listen("tcp", address)
	s.listener = lis
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return err
}

func  (s *RpcServer) Broadcast(command interface{}) error {
	var in pb.Messages
	switch v := command.(type) {
	case protocol.MessageCommand:
		in = pb.Messages{
			Name:                 v.Name,
			Message:              v.Message,
		}
	}
	for _, client := range s.clients {
		err := client.stream.Send(&in)
		if (err != nil) {
			return err
		}
 	}
 	return nil
}

func (s *RpcServer) Close() {}
