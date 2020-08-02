package common

import (
	"encoding/json"
	"fmt"
)

func Marshall(messageToMarshall Message) ([]byte, error) {
	fmt.Println("Marshall: marshalling "+messageToMarshall.Operation)
	return json.Marshal(messageToMarshall)
}

func Unmarshall(bytes []byte, bytesToMessage *Message) error {
	fmt.Println("Marshall: unmarshalling")
	return json.Unmarshal(bytes, bytesToMessage)
}
