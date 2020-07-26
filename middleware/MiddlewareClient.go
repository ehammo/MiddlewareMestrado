package middleware

import (
	infra "../infraestrutura"
	"fmt"
	"strings"
)

type MiddlewareClient struct {
	crh               *infra.ClientRequestHandler
	serverAddr        string
	currentLane       string
}

func NewMiddlewareClient(transportType string, currentLane string) *MiddlewareClient {
	return &MiddlewareClient {
		crh:               infra.NewClient(transportType),
		serverAddr:        "localhost:1111",
		currentLane:       currentLane,
	}
}

func (mc *MiddlewareClient) Register() {
	var msg = "REGISTER: "+mc.currentLane+"\n"
	mc.crh.Send([]byte(msg))
}

func  (mc *MiddlewareClient) ChangeLane(lane string) {
	var msg = "LANE: "+lane+"\n"
	mc.currentLane = lane
	mc.crh.Send([]byte(msg))
}

func (mc *MiddlewareClient) BroadcastMessage() {
	var msg = "BREAK: "+mc.currentLane+"\n"
	mc.crh.Send([]byte(msg))
}

func (mc *MiddlewareClient) BenchmarkMessages(qtd int) {
	for i := 0; i < qtd; i++ {
		mc.BroadcastMessage()
	}
}

func (mc *MiddlewareClient) Start() {
	mc.crh.Dial(mc.serverAddr)
	go mc.startLoop()
}
func (mc *MiddlewareClient) startLoop() {
	for {
		var data = mc.crh.Receive()
		var cmd = string(data)
		if strings.Contains(cmd, "BREAK") {
			// then the car should break
			fmt.Println("breaking car")
		} else {
			fmt.Println(cmd)
		}
	}
}