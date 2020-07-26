package middleware

import (
	infra "../infraestrutura"
	"fmt"
	"log"
	"strings"
	"sync"
)

type MiddlewareServer struct {
	srh               *infra.ServerRequestHandler
	serverAddr        string
	clients           []*client
	mutex             *sync.Mutex
}

type client struct {
	id int
	currentLane string
}

func NewMiddlewareServer(transportType string) *MiddlewareServer {
	return &MiddlewareServer{
		srh:               infra.NewServer(transportType),
		serverAddr:        "localhost:1111",
		mutex:             &sync.Mutex{},
		clients:           make([]*client, 10),
	}
}

func (ms *MiddlewareServer) Start() {
	ms.srh.Listen(ms.serverAddr)
	for {
		var id = ms.srh.AcceptNewClient()
		if id == -1 {
			break
		}
		var client = &client{
			id: id,
			currentLane: "UNKNOWN",
		}
		ms.mutex.Lock()
		ms.clients[id] = client
		ms.mutex.Unlock()
		go ms.serve(client)
	}
}

func (ms *MiddlewareServer) serve(c *client) {
	for {
		var data = ms.srh.Receive(c.id)
		log.Println(string(data))
		var cmd = strings.Split(string(data), ":")
		if cmd[0] == "REGISTER" || cmd[0] == "LANE" {
			c.currentLane = cmd[1]
			ms.mutex.Lock()
			ms.clients[c.id] = c
			ms.mutex.Unlock()
		} else if cmd[0] == "BREAK" {
			var lane = cmd[1]
			ms.mutex.Lock()
			for _, client := range ms.clients {
				if client != nil{
					fmt.Println(client.currentLane)
				}
				if client != nil && strings.Contains(lane, client.currentLane) {
					fmt.Println("sending this data from client ", len(data))
					fmt.Println("sending this string from client ", string(data))
					ms.srh.Send(data, client.id)
				}
			}
			ms.mutex.Unlock()
		}
	}
}
