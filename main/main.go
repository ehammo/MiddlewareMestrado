package main

import (
	myrpc "./rpc"
	client "./socket/client"
	server "./socket/server"
	"fmt"
	"log"
	"io"
	"math"
	"time"
)

func tcpServerStart() {
	var s server.ChatServer
	s = server.NewServer()
	err := s.Listen("192.168.0.16:1111")
	if err != nil {
		log.Printf("error=%s",err.Error())
	}
	go s.Start()
}

func UdpServerStart() {
	var s server.ChatServer
	s = server.NewUdpServer()
	err := s.Listen("192.168.0.16:2222")
	if err != nil {
		log.Printf("error=%s",err.Error())
	}
	go s.Start()
}

func rpcServerStart() {
	var s server.ChatServer
	s = myrpc.NewRpcServer()
	err := s.Listen("192.168.0.16:3333")
	if err != nil {
		log.Printf("error=%s",err.Error())
	}
	go s.Start()
}

func rpcClient() client.ChatClient {
	var c1 client.ChatClient
	c1 = myrpc.NewRpcClient()
	println("lets dial")
	err := c1.Dial("192.168.0.16:3333")
	if err != nil {
		log.Printf("main dial error=%s",err.Error())
	}
	go c1.Start()
	return c1
}

func udpClient() client.ChatClient {
	var c1 client.ChatClient
	c1 = client.NewUdpClient()
	err := c1.Dial("192.168.0.16:2222")
	if err != nil {
		log.Printf("main dial error=%s",err.Error())
	}
	go c1.Start()
	return c1
}

func tcpClient() client.ChatClient {
	var c1 client.ChatClient
	c1 = client.NewClient()
	c1.Dial("192.168.0.16:1111")
	go c1.Start()
	return c1
}

func runMessages(c1 client.ChatClient, sentMessages float64, currentSum float64, clientType string, shouldRead bool) {
	const defaultValue  = 10000
	log.Printf("sending %d messages",defaultValue)
	var sum, mean, sd, total float64
	total = defaultValue - sentMessages
	sum = currentSum
	var times [defaultValue]int64
	var forcefulBreak = false
	for i := 0; i < defaultValue; i++ {
		var t1 = time.Now()
		err := c1.SendMessage(fmt.Sprintf("%d", i))
		if err != nil {
			log.Printf("error=%s",err.Error())
			if err == io.EOF {
				var NewClient = createClient(clientType)
				go runMessages(NewClient, float64(i), sum, clientType, shouldRead)
				i = defaultValue
				forcefulBreak = true
			}
		} else {
			time.Sleep(30*time.Millisecond)
			var t2 = time.Now()
			var delay = 30*time.Millisecond
			t1 = t1.Add(delay)
			times[i] = t2.Sub(t1).Nanoseconds()
			if float64(times[i]) == 0 {
				total -= 1
			}
			sum += float64(times[i])	
		}
	}
	println("sent")
	if (shouldRead) {
		println("reading")
		var t1 = time.Now()
		count := 0
		for i := range c1.Incoming() {
			println(count, i.Message)
			count+=1
		}
		var t2 = time.Now()	
		sum += float64(t2.Sub(t1).Nanoseconds())
	}
	println("read")
	if (!forcefulBreak && shouldRead) {
		if sum > 0 && total>0 {
			mean = sum/total
		} else {
			mean = 0
		}
	
		log.Printf("Mean: %f", mean)
		sd = 0
		for i := 0; i < defaultValue; i++ {
			sd += math.Pow(float64(times[i]) - mean, 2)
		}
		sd = math.Sqrt(sd/defaultValue)
		log.Printf("Sd: %f", sd)
	} else {
		println("forceful break")
	}
	println("finishing runmessages")
	c1.Clean()
}

func createClient(clientType string) client.ChatClient {
	if ("rpc" == clientType) {
		return rpcClient()
	} else if ("udp" == clientType) {
		return udpClient()
	} else {
		return tcpClient()
	}
}

func runFiveClients(clientType string) {
	var shouldRead, shouldReadC1 bool
	shouldReadC1 = true
	if (clientType != "udp") {
		shouldRead = true
	} else {
		shouldRead = false
	}
	// log.Printf("Com 1 cliente")
	// log.Printf(clientType)
	// var c1 = createClient(clientType)
	// c1.SetName("c1")
	// time.Sleep(1*time.Second)
	// go runMessages(c1, 0, 0, clientType, shouldReadC1)
	// fmt.Scanln()
	// log.Printf("Closing client")
	// c1.Close()
	// fmt.Scanln()
	
	// log.Printf("Com 2 clientes")
	// log.Printf(clientType)
	// c1 = createClient(clientType)
	// var c2 = createClient(clientType)
	// c1.SetName("c1")
	// c2.SetName("c2")
	// time.Sleep(1*time.Second)
	// go runMessages(c2, 0, 0, clientType, shouldRead)
	// go runMessages(c1, 0, 0, clientType, shouldReadC1)
	// fmt.Scanln()
	// log.Printf("Closing client")
	// c1.Close()
	// c2.Close()
	// fmt.Scanln()

	log.Printf("Com 3 clientes")
	log.Printf(clientType)
	var c1 = createClient(clientType)
	var c2 = createClient(clientType)
	var c3 = createClient(clientType)
	c1.SetName("c1")
	c2.SetName("c2")
	c3.SetName("c3")
	time.Sleep(1*time.Second)
	go runMessages(c1, 0, 0, clientType, shouldReadC1)
	go runMessages(c2, 0, 0, clientType, shouldRead)
	go runMessages(c3, 0, 0, clientType, shouldRead)
	fmt.Scanln()
	log.Printf("Closing client")
	c1.Close()
	c2.Close()
	c3.Close()
	fmt.Scanln()

	log.Printf("Com 4 clientes")
	log.Printf(clientType)
	c1 = createClient(clientType)
	c2 = createClient(clientType)
	c3 = createClient(clientType)
	var c4 = createClient(clientType)
	c1.SetName("c1")
	c2.SetName("c2")
	c3.SetName("c3")
	c4.SetName("c4")
	time.Sleep(1*time.Second)
	go runMessages(c1, 0, 0, clientType, shouldReadC1)
	go runMessages(c2, 0, 0, clientType, shouldRead)
	go runMessages(c3, 0, 0, clientType, shouldRead)
	go runMessages(c4, 0, 0, clientType, shouldRead)
	fmt.Scanln()
	log.Printf("Closing client")
	c1.Close()
	c2.Close()
	c3.Close()
	c4.Close()
	fmt.Scanln()

	log.Printf("Com 5 clientes")
	log.Printf(clientType)
	c1 = createClient(clientType)
	c2 = createClient(clientType)
	c3 = createClient(clientType)
	c4 = createClient(clientType)
	var c5 = createClient(clientType)
	c1.SetName("c1")
	c2.SetName("c2")
	c3.SetName("c3")
	c4.SetName("c4")
	c5.SetName("c5")
	time.Sleep(1*time.Second)
	go runMessages(c1, 0, 0, clientType, shouldReadC1)
	go runMessages(c2, 0, 0, clientType, shouldRead)
	go runMessages(c3, 0, 0, clientType, shouldRead)
	go runMessages(c4, 0, 0, clientType, shouldRead)
	go runMessages(c5, 0, 0, clientType, shouldRead)
	fmt.Scanln()
	defer c1.Close()
	defer c2.Close()
	defer c3.Close()
	defer c4.Close()
	defer c5.Close()
}

func main() {
	// runFiveClients("rpc")
	// runFiveClients("tcp")
	runFiveClients("udp")
	

}
