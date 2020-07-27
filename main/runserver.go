package main

import (
	"fmt"
	middleware "../middleware"
)


func main() {
	var s = middleware.NewMiddlewareServer("tcp")
	go s.Start()
	fmt.Scanln()
	// var s = middleware.NewMiddlewareServer("udp")
	// go s.Start()
	// fmt.Scanln()
}
