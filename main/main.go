package main

import (
	myrpc "./rpc"
	client "./socket/client"
	server "./socket/server"
	"fmt"
	"io"
	"log"
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
	var err = c1.Dial("192.168.0.16:1111")
	if err != nil {
		log.Printf("main dial error=%s",err.Error())
	}
	go c1.Start()
	return c1
}

func runMessages(c1 client.ChatClient, sentMessages int, currentSum float64, clientType string, shouldRead bool,
	             times [10000]float64) {
	const defaultValue  = 10000
	log.Printf("sending %d messages",defaultValue)
	var sum float64
	var total int
	total = defaultValue - sentMessages
	sum = currentSum
	var forcefulBreak = false
	var delay = 1*time.Millisecond
	for i := 0; i < defaultValue; i++ {
		var t1 = time.Now()
		err := c1.SendMessage(fmt.Sprintf("%d", i))
		if err != nil {
			log.Printf("error=%s",err.Error())
			if err == io.EOF {
				var NewClient = createClient(clientType)
				go runMessages(NewClient, i, sum, clientType, shouldRead, times)
				i = defaultValue
				forcefulBreak = true
			}
		} else {
			time.Sleep(delay)
			t1 = t1.Add(delay)
			times[i] = float64(time.Since(t1).Nanoseconds())
			if times[i] == 0 {
				println("Eitcha deu um 0")
				total -= 1
			}
			sum += times[i]
		}
	}
	println("sent")
	if shouldRead {
		println("reading")
		var t1 = time.Now()
		count := 0
		for i := range c1.Incoming() {
			println(count, i.Message)
			count+=1
		}
		time.Sleep(delay)
		t1 = t1.Add(delay)
		var tReadFinal = float64(time.Since(t1).Nanoseconds())
		if tReadFinal == 0 {
			println("Eitcha deu um 0 na leitura")
		}
		sum += float64(time.Since(t1).Nanoseconds())
	}
	println("read")
	if !forcefulBreak && shouldRead {

	} else {
		println("forceful break")
	}
	println("finishing runmessages")
	c1.Clean()
}

func calculateMeanAndSd(total int, times [] float64, sum float64) {
	var mean, sd float64
	if sum > 0 && total > 0 {
		mean = sum/float64(total)
	} else {
		mean = 0
	}

	log.Printf("Mean: %f", mean)
	sd = 0

	for i := 0; i < total; i++ {
		if times[i] > 0 {
			sd += math.Pow(times[i] - mean, 2)
		}
	}
	sd = math.Sqrt(sd/float64(total))
	log.Printf("Sd: %f", sd)
}

func createClient(clientType string) client.ChatClient {
	if "rpc" == clientType {
		return rpcClient()
	} else if "udp" == clientType {
		return udpClient()
	} else {
		return tcpClient()
	}
}

func runFiveClients(clientType string) {
	var shouldRead bool
	if clientType != "udp" {
		shouldRead = true
	} else {
		shouldRead = false
	}
	var times [10000]float64


	log.Printf("Com 5 clientes")
	log.Printf(clientType)
	var c1 = createClient(clientType)
	var c2 = createClient(clientType)
	var c3 = createClient(clientType)
	var c4 = createClient(clientType)
	var c5 = createClient(clientType)
	c1.SetName("c1")
	c2.SetName("c2")
	c3.SetName("c3")
	c4.SetName("c4")
	c5.SetName("c5")
	time.Sleep(1*time.Second)
	go runMessages(c1, 0, 0, clientType, shouldRead,times)
	go runMessages(c2, 0, 0, clientType, shouldRead,times)
	go runMessages(c3, 0, 0, clientType, shouldRead,times)
	go runMessages(c4, 0, 0, clientType, shouldRead,times)
	go runMessages(c5, 0, 0, clientType, shouldRead,times)
	fmt.Scanln()
	defer c1.Close()
	defer c2.Close()
	defer c3.Close()
	defer c4.Close()
	defer c5.Close()
}

func main() {
	runFiveClients("rpc")
	// runFiveClients("tcp")
	// runFiveClients("udp")
}
