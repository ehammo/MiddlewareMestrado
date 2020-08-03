package common

type Message struct {
	Operation string `json:"operation"`
	Topic string `json:"topic"`
}

type Invocation struct {
	Addr string
	Message *Message
}

type Termination struct {
	Result string
}

type AOR struct {
	address string
	protocol string
	objectId string
}

func (aor *AOR) equals(aor2 *AOR) bool {
	return aor.objectId == aor2.objectId &&
		aor.address == aor2.address &&
		aor.protocol == aor2.protocol
}

func NewAOR(address string, protocol string, objectId string) *AOR {
	return &AOR{
		address:  address,
		protocol: protocol,
		objectId: objectId,
	}
}
