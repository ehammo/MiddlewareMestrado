package main

import (
	"fmt"
	"log"
	"time"

	d "../distribution"
	n "../naming"
)

func startClient(id int) *d.ClientProxy {
	namingProxy := n.NewNamingProxy("localhost:1243")

	aor := namingProxy.LookUp("Vanet")
	fmt.Println("Received aor:")
	fmt.Println(aor)
	var c = d.NewClientProxy(aor, id)
	go c.Start()
	return c
}

func threeBreakingCars() {
	log.Printf("twoBreakingCars")
	var c1, c2, c3, c4, c5 *d.ClientProxy
	c1 = startClient(0)
	c2 = startClient(1)
	c3 = startClient(2)
	c4 = startClient(3)
	c5 = startClient(4)
	c1.RegisterKey()
	c1.RegisterOnLane("lane1")
	c2.RegisterKey()
	c3.RegisterKey()
	c4.RegisterKey()
	c5.RegisterKey()
	c2.RegisterOnLane("lane1")
	c3.RegisterOnLane("lane2")
	c4.RegisterOnLane("lane2")
	c5.RegisterOnLane("lane2")
	time.Sleep(5 * time.Second)
	fmt.Println("5 seconds to go")
	time.Sleep(5 * time.Second)
	c1.BroadcastEvent("lane2")
}

func main() {
	threeBreakingCars()
	fmt.Scanln()
	// twoBreakingCars()
	// fmt.Scanln()
}
