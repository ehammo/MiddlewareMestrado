package distribution

import (
	common "../common"
	"fmt"
	"strings"
	"sync"
)

type EventBus struct {
	client  *Client
	invoker *Invoker
	mutex   *sync.Mutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		client: nil,
		invoker: nil,
		mutex: &sync.Mutex{},
	}
}

func (e *EventBus) ChangeLane(newLane string) string {
	return e.handleEvent("CHANGE", newLane)
}

func (e *EventBus) BroadcastEvent(lane string) string {
	return e.handleEvent("BREAK", lane)
}

func (e *EventBus) RegisterOnLane(lane string) string {
	return e.handleEvent("REGISTER", lane)
}

func (e *EventBus) handleEvent(op string, lane string) string {
	fmt.Println("Calling handle event from eventbus")
	if e.client == nil || e.invoker == nil {
		return "Error nil invoker or nil client"
	}
	e.mutex.Lock()
	e.invoker.mutex.Lock()
	if op == "REGISTER" || op == "LANE" {
		e.client.currentLane = lane
		e.invoker.clients[e.client.id] = e.client
	} else if op == "BREAK" {
		for _, client := range e.invoker.clients {
			if client != nil && strings.Contains(lane, client.currentLane) {
				message := &common.Message{
					Operation: op,
					Topic:     lane,
				}
				e.invoker.sendMessage(message, client)
			}
		}
	} else {
		return "Invalid operation"
	}
	e.invoker.mutex.Unlock()
	e.mutex.Unlock()
	return "Success"
}

func (e *EventBus) SetClient(c *Client) {
	e.mutex.Lock()
	e.client = c
	e.mutex.Unlock()
}

func (e *EventBus) SetInvoker(i *Invoker) {
	e.mutex.Lock()
	e.invoker = i
	e.mutex.Unlock()
}
