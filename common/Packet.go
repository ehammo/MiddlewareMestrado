package common

import (
	"encoding/json"
	"fmt"
)

type Packet struct {
	Header []byte
	Body   []byte
}

func NewRequestPacket(message Message) *Packet {
	op := []byte(message.Operation)
	topic := []byte(message.Topic)
	return &Packet{
		Header: op,
		Body:   topic,
	}
}

func NewReplyPacket(response string) *Packet {
	header := []byte("Reply")
	body := []byte(response)
	return &Packet{
		Header: header,
		Body:   body,
	}
}

func NewLookUpReplyPacket(aor interface{}) *Packet {
	fmt.Println("Creating reply package")
	header := []byte("lookup")
	body, _ := json.Marshal(aor)
	return &Packet{
		Header: header,
		Body:   body,
	}
}

func NewLookUpRequestPacket(message Message, aor interface{}) *Packet {
	header := []byte("lookup")
	aorBody, _ := json.Marshal(aor)
	messageBody, _ := json.Marshal(message)
	divider := make([]byte, 2)
	divider[0] = '\n'
	divider[1] = '\n'
	aorBodyDivider := append(aorBody, divider...)
	body := append(aorBodyDivider, messageBody...)
	return &Packet{
		Header: header,
		Body:   body,
	}
}

func CreateLookupMessageFromLookupPacket(packet *Packet) *LookupMessage {
	if string(packet.Header) == "lookup" {
		var aorBody []byte
		var messageBody []byte
		var lastOne = false
		for i, b := range packet.Body {
			if b == '\n' {
				if lastOne == false {
					lastOne = true
				} else {
					messageBody = packet.Body[i:len(packet.Body)]
					aorBody = packet.Body[0:i]
					break
				}
			}
		}
		aor := &AOR{}
		message := &Message{}
		_ = json.Unmarshal(aorBody, aor)
		_ = json.Unmarshal(messageBody, message)
		return &LookupMessage{
			Message: message,
			AOR:    aor,
		}
	}
	return nil
}
