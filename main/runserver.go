package main

import (
	"fmt"
	d "../distribution"
)


func main() {
	var s = d.NewInvoker("localhost:1111", "tcp")
	go s.Start()
	fmt.Scanln()
	var s2 = d.NewInvoker("localhost:1111", "udp")
	go s2.Start()
	fmt.Scanln()
}
