package main

import (
	"fmt"

	c "../common"
)

func main() {
	srvkp := c.GenerateKeypair(false)
	kp := c.GenerateKeypair(true)

	pub := &c.AORPUB{
		N: (*(kp.Pub.N)).String(),
		E: kp.Pub.E,
	}

	srvkppub := &c.AORPUB{
		N: (*(srvkp.Pub.N)).String(),
		E: srvkp.Pub.E,
	}

	fmt.Println(len(srvkppub.N))

	var message = &c.Message{
		Operation: "RegisterKey",
		Topic:     pub,
	}
	packet := c.NewRequestPacket(message)
	//fmt.Println(packet)
	marshalledMessage, _ := c.Marshall(*packet)
	criptedMessage := c.Encrypt(marshalledMessage, srvkp.Pub, false)
	decrypted := c.Decrypt(criptedMessage, srvkp.Priv, false)
	c.Unmarshall(decrypted, packet)
	//decryptedM := packet.Body.ReqBody.Body[0].(*c.AORPUB)
	//fmt.Println(decryptedM)

	message = &c.Message{
		Operation: "BroadcastMessage",
		Topic:     "lane1",
	}
	packet = c.NewRequestPacket(message)
	fmt.Println(packet)
	marshalledMessage, _ = c.Marshall(*packet)

	criptedMessage = c.Encrypt(marshalledMessage, kp.Pub, true)
	decrypted = c.Decrypt(criptedMessage, kp.Priv, true)
	c.Unmarshall(decrypted, packet)
	decrypted2 := packet.Body.ReqBody.Body[0].(string)
	fmt.Println(decrypted2)

	//  criptedMessage = c.Encrypt(marshalledMessage, kp.Pub, true)
	//  decrypted = c.Decrypt(criptedMessage, kp.Priv, true)
	//	c.Unmarshall(decrypted, packet)
	//	decrypted2 = packet.Body.ReqBody.Body[0].(string)
	//	fmt.Println(decrypted2)
}
