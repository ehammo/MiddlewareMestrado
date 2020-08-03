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
	Address string
	Protocol string
	ObjectId string
}

func (aor *AOR) equals(aor2 *AOR) bool {
	return aor.ObjectId == aor2.ObjectId &&
		aor.Address == aor2.Address &&
		aor.Protocol == aor2.Protocol
}

func NewAOR(address string, protocol string, objectId string) *AOR {
	return &AOR{
		Address:  address,
		Protocol: protocol,
		ObjectId: objectId,
	}
}
