package main

import (
	"fmt"
	"log"

	d "../distribution"
	n "../naming"
)

func startClient() *d.ClientProxy {
	namingProxy := n.NewNamingProxy("localhost:1243")

	aor := namingProxy.LookUp("Vanet")
	fmt.Println("Received aor:")
	fmt.Println(aor)
	var c = d.NewClientProxy(aor)
	//c.Start()
	return c
}

//     c3   c2  c1^  (*)
//     c5   c4
func threeBreakingCars() {
	log.Printf("twoBreakingCars")
	var c1, c2, c3, c4, c5 *d.ClientProxy
	c1 = startClient()
	c2 = startClient()
	c3 = startClient()
	c4 = startClient()
	c5 = startClient()
	c1.RegisterKey()
	c2.RegisterKey()
	c3.RegisterKey()
	c4.RegisterKey()
	c5.RegisterKey()
	c1.RegisterOnLane("lane1")
	c2.RegisterOnLane("lane1")
	c3.RegisterOnLane("lane2")
	c4.RegisterOnLane("lane2")
	c5.RegisterOnLane("lane2")
	c1.BroadcastEvent("lane2")
}

func main() {
	threeBreakingCars()
	fmt.Scanln()
	// twoBreakingCars()
	// fmt.Scanln()
}
