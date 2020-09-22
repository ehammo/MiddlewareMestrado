package main

import (
	"fmt"

	naming "../naming"
)

func main() {
	invoker := naming.NewNamingInvoker("172.17.0.2:1243")
	invoker.Start()
	fmt.Scanln()
}
