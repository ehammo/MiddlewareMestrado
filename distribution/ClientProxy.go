package distribution

import (
	"crypto/rsa"
	"fmt"
	"math/big"

	common "../common"
	infra "../infra"
)

type ClientProxy struct {
	Kp         *common.Keypair
	protocol   string
	srvAddress string
	id         string
	ClientId   int
	crh        *infra.ClientRequestHandler
	srvPub     *rsa.PublicKey
}

func NewClientProxy(aor *common.AOR, id int) *ClientProxy {
	N := big.Int{}
	//N.SetString(aor.Pub.N, 10)
	N.SetString(aor.N, 10)
	srvPub := &rsa.PublicKey{
		N: &N,
		E: aor.E,
		//E: aor.Pub.E,
	}
	return &ClientProxy{
		protocol:   aor.Protocol,
		srvAddress: aor.Address,
		id:         aor.ObjectId,
		Kp:         common.GenerateKeypair(true),
		srvPub:     srvPub,
		ClientId:   id,
	}
}

func (cp *ClientProxy) ChangeLane(newLane string) {
	cp.invokeCommand("CHANGE", newLane)
}

func (cp *ClientProxy) BroadcastEvent(lane string) {
	cp.invokeCommand("BroadcastEvent", lane)
}

func (cp *ClientProxy) RegisterOnLane(lane string) {
	cp.invokeCommand("Register", lane)
}

func (cp *ClientProxy) RegisterKey() {
	pub := &common.AORPUB{
		N: (*(cp.Kp.Pub.N)).String(),
		E: cp.Kp.Pub.E,
	}
	cp.invokeCommand("RegisterKey", pub)
}

func (cp *ClientProxy) invokeCommand(op string, object interface{}) {
	fmt.Println("ClientProxy: invoking " + op)
	var message = &common.Message{
		Operation: op,
		Topic:     object,
		ClientId:  cp.ClientId,
	}
	var invocation = &common.Invocation{
		Addr:    cp.srvAddress,
		Message: message,
	}
	result := cp.Invoke(invocation).Result
	message, ok := result.(*common.Message)
	if !ok {
		fmt.Println(result.(string))
	} else {
		fmt.Println(message)
	}
}

func (cp *ClientProxy) Invoke(invocation *common.Invocation) *common.Termination {
	message := invocation.Message
	if cp.crh == nil {
		cp.crh = infra.NewCRH(invocation.Addr)
	}
	packet := common.NewRequestPacket(message)
	fmt.Println(packet)
	marshalledMessage, err := common.Marshall(*packet)

	if err != nil {
		return marshallingError(err)
	}
	fmt.Println("Encrypting")
	criptedMessage := common.Encrypt(marshalledMessage, cp.srvPub, false)
	fmt.Println("Encrypted")
	if cp.protocol == "tcp" {
		return cp.SendReceiveTcp(criptedMessage, message.IsReplyRequired())
	} else {
		return cp.SendReceiveUDP(criptedMessage, message.IsReplyRequired())
	}
}

func (cp *ClientProxy) SendReceiveTcp(message []byte, isReplyRequired bool) *common.Termination {
	fmt.Println("Sending")
	cp.crh.SendTcp(message)
	fmt.Println("Sent")
	return &common.Termination{Result: &common.Message{Topic: "Success"}}
}

func unpack(data []byte) *common.Termination {
	packet := &common.Packet{}
	err := common.Unmarshall(data, packet)
	if err != nil {
		return marshallingError(err)
	}
	var ter *common.Termination
	if len(packet.Body.RepBody.Body) > 0 {
		ter = &common.Termination{Result: packet.Body.RepBody.Body[0]}
	} else {
		ter = &common.Termination{Result: "No response"}
	}

	return ter
}

func (cp *ClientProxy) SendReceiveUDP(message []byte, isReplyRequired bool) *common.Termination {
	cp.crh.SendUDP(message)
	if isReplyRequired {
		data := cp.crh.ReceiveUDP()
		decriptedData := common.Decrypt(data, cp.Kp.Priv, true)
		return unpack(decriptedData)
	}
	return &common.Termination{Result: &common.Message{Topic: "Success"}}
}

func (cp *ClientProxy) Start() {
	for {
		var data []byte
		if cp.crh != nil {
			if cp.protocol == "tcp" {
				data = cp.crh.ReceiveTcp()
			} else {
				data = cp.crh.ReceiveUDP()
			}
			decryptedData := common.Decrypt(data, cp.Kp.Priv, true)
			ter := unpack(decryptedData)
			result := ter.Result
			fmt.Println("Start-result")
			if result != nil {
				fmt.Println(result)
			}
		}
	}
}
