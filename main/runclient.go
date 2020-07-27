package main

import (
	middleware "../middleware"
	"fmt"
	"log"
)

func startClient(clientType string, lane string) *middleware.MiddlewareClient {
	var c = middleware.NewMiddlewareClient(clientType, lane)
	c.Start()
	return c
}

//     c3   c2  c1^  (*)
//     c5   c4
func threeBreakingCars(clientType string) {
	log.Printf("twoBreakingCars")
	log.Printf(clientType)
	var c1,c2,c3,c4,c5 *middleware.MiddlewareClient
	c1 = startClient(clientType, "lane1")
	c2 = startClient(clientType, "lane1")
	c3 = startClient(clientType, "lane1")
	c4 = startClient(clientType, "lane2")
	c5 = startClient(clientType, "lane2")
	c1.Register()
	c2.Register()
	c3.Register()
	c4.Register()
	c5.Register()
	c1.BroadcastMessage()
}

//     c3   c2  c1^
//     c5   c4^  (*)
func twoBreakingCars(clientType string) {
	log.Printf("twoBreakingCars")
	log.Printf(clientType)
	var c1,c2,c3,c4,c5 *middleware.MiddlewareClient
	c1 = startClient(clientType, "lane1")
	c2 = startClient(clientType, "lane1")
	c3 = startClient(clientType, "lane1")
	c4 = startClient(clientType, "lane2")
	c5 = startClient(clientType, "lane2")
	c1.Register()
	c2.Register()
	c3.Register()
	c4.Register()
	c5.Register()
	c4.BroadcastMessage()
}

func main() {
	threeBreakingCars("tcp")
	fmt.Scanln()
	// twoBreakingCars("udp")
	// fmt.Scanln()
}
