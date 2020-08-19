package main

import (
	c "../common"
	d "../distribution"
	n "../naming"
	"fmt"
)


func main() {
	namingProxy := n.NewNamingProxy()
	aor := &c.AOR{
		Address:  "localhost:1111",
		Protocol: "tcp",
		ObjectId: "1",
	}
	namingProxy.Register("Vanet", aor)
	var s = d.NewInvoker(aor.Address, aor.Protocol)
	go s.Start()
	fmt.Scanln()
}
