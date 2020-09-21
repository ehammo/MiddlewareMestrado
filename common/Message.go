package common

import "strconv"

type Message struct {
	Operation string      `json:"operation"`
	Topic     interface{} `json:"topic"`
	AOR       interface{} `json:"AOR"`
}

func (m *Message) IsReplyRequired() bool {
	if m.Operation == "LookUp" || m.Operation == "Register" {
		return true
	}
	return false
}

type Invocation struct {
	Addr    string
	Message *Message
}

type Termination struct {
	Result interface{}
}

type AOR struct {
	Address  string `json:"Address"`
	Protocol string `json:"Protocol"`
	ObjectId string `json:"ObjectId"`
	N        string `json:"N"`
	E        int    `json:"E"`
	//	Pub      AORPUB
}

type AORPUB struct {
	N string
	E int
}

func (aor *AOR) equals(aor2 *AOR) bool {
	return aor.ObjectId == aor2.ObjectId &&
		aor.Address == aor2.Address &&
		aor.Protocol == aor2.Protocol &&
		//	aor.Pub.N == aor2.Pub.N &&
		//	aor.Pub.E == aor2.Pub.E
		aor.N == aor2.N &&
		aor.E == aor2.E
}

func (aor *AOR) ToString() string {
	return aor.Address + " " + aor.Protocol + " " + aor.N + " " + strconv.Itoa(aor.E)
	//	return aor.Address + " " + aor.Protocol + " " + aor.Pub.N + " " + strconv.Itoa(aor.Pub.E)
}
