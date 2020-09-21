package distribution

import (
	"fmt"
	"strings"
	"sync"

	common "../common"
)

type EventBus struct {
	client       *Client
	queueManager *QueueManager
	mutex        *sync.Mutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		client:       nil,
		queueManager: nil,
		mutex:        &sync.Mutex{},
	}
}

func (e *EventBus) ChangeLane(newLane string) string {
	return e.handleEvent("CHANGE", newLane)
}

func (e *EventBus) BroadcastEvent(lane string) string {
	return e.handleEvent("BroadcastEvent", lane)
}

func (e *EventBus) RegisterOnLane(lane string) string {
	return e.handleEvent("Register", lane)
}

func (e *EventBus) handleEvent(op string, lane string) string {
	fmt.Println("Calling handle event from eventbus")
	if e.client == nil || e.queueManager == nil {
		return "Error nil queueManager or nil client"
	}
	e.mutex.Lock()
	fmt.Println("Lock1 ok")
	e.queueManager.mutex.Lock()
	fmt.Println("Lock2 ok")
	if op == "Register" {
		e.client.currentLane = lane
		e.queueManager.clients[e.client.Id] = e.client
		message := &common.Message{
			Operation: op,
			Topic:     e.client.Id,
		}
		e.queueManager.sendMessage(message, e.client)
	} else if op == "ChangeLane" {
		e.client.currentLane = lane
		e.queueManager.clients[e.client.Id] = e.client
	} else if op == "BroadcastEvent" {
		for _, client := range e.queueManager.clients {
			if client != nil && strings.Contains(lane, client.currentLane) {
				message := &common.Message{
					Operation: "break",
					Topic:     "",
				}
				e.queueManager.sendMessage(message, client)
			}
		}
	} else {
		return "Invalid operation"
	}
	e.queueManager.mutex.Unlock()
	e.mutex.Unlock()
	return "Success"
}

func (e *EventBus) SetClient(c *Client) {
	e.mutex.Lock()
	e.client = c
	e.mutex.Unlock()
}

func (e *EventBus) SetQueueManager(qm *QueueManager) {
	e.mutex.Lock()
	e.queueManager = qm
	e.mutex.Unlock()
}
