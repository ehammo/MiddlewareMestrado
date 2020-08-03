package common

import (
	"encoding/json"
)

func Marshall(messageToMarshall Packet) ([]byte, error) {
	return json.Marshal(messageToMarshall)
}

func Unmarshall(bytes []byte, bytesToMessage *Packet) error {
	return json.Unmarshal(bytes, bytesToMessage)
}
