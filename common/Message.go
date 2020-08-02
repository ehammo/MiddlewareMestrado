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
	AOR string
}

func NewAOR(address string, protocol string, objectId string) *AOR {
	return &AOR{
		AOR: address+protocol+objectId,
	}
}
