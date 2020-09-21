package common

import "fmt"

type Packet struct {
	Header Header
	Body   Body
}

type Header struct {
	Version  string
	ClientId int
}

type Body struct {
	ReqHeader ReqHeader
	ReqBody   ReqRepBody
	RepHeader RepHeader
	RepBody   ReqRepBody
}

type ReqHeader struct {
	ResponseExpected bool
	Operation        string
}

type ReqRepBody struct {
	Body []interface{}
}

type RepHeader struct {
	status string
}

func NewRequestPacket(message *Message) *Packet {
	fmt.Println("Creating package")
	reqHeader := &ReqHeader{
		ResponseExpected: message.IsReplyRequired(),
		Operation:        message.Operation,
	}
	var reqReqBodyArray = make([]interface{}, 2)
	reqReqBodyArray[0] = message.Topic
	reqReqBodyArray[1] = message.AOR
	reqRepBody := &ReqRepBody{
		Body: reqReqBodyArray,
	}
	packet := &Packet{
		Header: Header{
			Version:  "1.0",
			ClientId: message.ClientId,
		},
		Body: Body{
			ReqHeader: *reqHeader,
			ReqBody:   *reqRepBody,
		},
	}
	return packet
}

func NewReplyPacket(response interface{}, status string) *Packet {
	repHeader := &RepHeader{status: status}
	var reqReqBodyArray = make([]interface{}, 1)
	reqReqBodyArray[0] = response
	reqRepBody := &ReqRepBody{
		Body: reqReqBodyArray,
	}
	return &Packet{
		Header: Header{
			Version: "1.0",
		},
		Body: Body{
			RepHeader: *repHeader,
			RepBody:   *reqRepBody,
		},
	}
}
