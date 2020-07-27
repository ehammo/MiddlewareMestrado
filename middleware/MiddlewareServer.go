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
	clients           []*infra.Client
	mutex             *sync.Mutex
	maxClients        int
	lastAddClient     int
}

func NewMiddlewareServer(transportType string) *MiddlewareServer {
	return &MiddlewareServer{
		srh:               infra.NewServer(transportType, "localhost:1111"),
		serverAddr:        "localhost:1111",
		mutex:             &sync.Mutex{},
		clients:           make([]*infra.Client, 10),
		maxClients:        10,
		lastAddClient:      0,
	}
}

func (ms *MiddlewareServer) Start() {
	for {
		var client = ms.srh.AcceptNewClient()
		if client.UniqueId == ms.maxClients {
			break
		}
		ms.mutex.Lock()
		ms.clients[client.UniqueId] = client
		ms.lastAddClient = client.UniqueId
		ms.mutex.Unlock()
		go ms.serve(client)
	}
}	

func (ms *MiddlewareServer) serve(c *infra.Client) {
	for {
		var data, addr = ms.srh.Receive(c)
		if c.Addr == nil {
			c.Addr = addr
		} else if c.Addr != addr {
			ms.mutex.Lock()
			var newclient = true
			for _, client := range ms.clients {
				if c.Addr == addr {
					newclient = false
					c = client
				}
			}
			if newclient {
				ms.lastAddClient = ms.lastAddClient+1
				newclientobject := c
				newclientobject.Addr = addr
				newclientobject.UniqueId = ms.lastAddClient
				ms.clients[ms.lastAddClient] = newclientobject
			}
			ms.mutex.Unlock()
		}
		log.Println(string(data))
		var cmd = strings.Split(string(data), ":")
		if cmd[0] == "REGISTER" || cmd[0] == "LANE" {
			c.CurrentLane = cmd[1]
			ms.mutex.Lock()
			ms.clients[c.UniqueId] = c
			ms.mutex.Unlock()
		} else if cmd[0] == "BREAK" {
			var lane = cmd[1]
			ms.mutex.Lock()
			for _, client := range ms.clients {
				if client != nil{
					fmt.Println(client.CurrentLane)
				}
				if client != nil && strings.Contains(lane, client.CurrentLane) {
					ms.srh.Send(data, client)
				}
			}
			ms.mutex.Unlock()
		}
	}
}
