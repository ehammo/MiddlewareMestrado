package main

import (
	common "../common"
	d "../distribution"
	"fmt"
	"log"
)

func startClient(clientType string, id string) *d.ClientProxy {
	aor := common.NewAOR("localhost:1111", clientType, id)
	var c = d.NewClientProxy(aor)
	//c.Start()
	return c
}

//     c3   c2  c1^  (*)
//     c5   c4
func threeBreakingCars(clientType string) {
	log.Printf("twoBreakingCars")
	log.Printf(clientType)
	var c1,c2,c3,c4,c5 *d.ClientProxy
	c1 = startClient(clientType, "vanetqueue")
	c2 = startClient(clientType, "vanetqueue")
	c3 = startClient(clientType, "vanetqueue")
	c4 = startClient(clientType, "vanetqueue")
	c5 = startClient(clientType, "vanetqueue")
	c1.RegisterOnLane("lane1")
	c2.RegisterOnLane("lane1")
	c3.RegisterOnLane("lane2")
	c4.RegisterOnLane("lane2")
	c5.RegisterOnLane("lane2")
	c1.BroadcastEvent("lane2")
}


func main() {
	threeBreakingCars("tcp")
	fmt.Scanln()
	// twoBreakingCars("udp")
	// fmt.Scanln()
}
