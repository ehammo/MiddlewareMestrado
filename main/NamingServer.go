package main

import (
	naming "../naming"
	"fmt"
)

func main() {
	invoker := naming.NewNamingInvoker("localhost:1243")
	invoker.Start()
	fmt.Scanln()
}