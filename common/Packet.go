package common

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
