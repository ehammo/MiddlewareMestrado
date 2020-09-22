package main

import (
	"fmt"

	c "../common"
	d "../distribution"
	n "../naming"
)

func main() {
	kp := c.GenerateKeypair(false)
	namingProxy := n.NewNamingProxy("172.17.0.2:1243")

	aor := &c.AOR{
		Address:  "172.17.0.3:1111",
		Protocol: "tcp",
		ObjectId: "1",
		N:        (*(kp.Pub.N)).String(),
		E:        kp.Pub.E,
	}
	namingProxy.Register("Vanet", aor)
	var s = d.NewQueueManager(aor.Address, aor.Protocol, kp)
	go s.Start()
	fmt.Scanln()
}
