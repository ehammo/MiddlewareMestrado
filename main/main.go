package main

import (
	rabbit "./rabbitmq"
	myrpc "./rpc"
	client "./socket/client"
	server "./socket/server"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

func tcpServerStart() {
	var s server.ChatServer
	s = server.NewServer()
	err := s.Listen("192.168.56.101:1111")
	if err != nil {
		log.Printf("error=%s",err.Error())
	}
	go s.Start()
}

func UdpServerStart() {
	var s server.ChatServer
	s = server.NewUdpServer()
	err := s.Listen("192.168.56.101:2222")
	if err != nil {
		log.Printf("error=%s",err.Error())
	}
	go s.Start()
}

func rpcServerStart() {
	var s server.ChatServer
	s = myrpc.NewRpcServer()
	err := s.Listen("192.168.56.101:3333")
	if err != nil {
		log.Printf("error=%s",err.Error())
	}
	go s.Start()
}

func rabbitServerStart(clients []string) {
	var s = rabbit.NewServer()
	for i := range clients {
		s.Register(clients[i])
	}
	go s.Start()
}

func rabbitClient(name string) client.ChatClient{
	var c1 client.ChatClient
	c1 = rabbit.NewClient()
	println("lets dial")
	err := c1.Dial(name)
	if err != nil {
		log.Printf("main dial error=%s",err.Error())
	}
	go c1.Start()
	return c1
}

func rpcClient() client.ChatClient {
	var c1 client.ChatClient
	c1 = myrpc.NewRpcClient()
	println("lets dial")
	err := c1.Dial("192.168.56.101:3333")
	if err != nil {
		log.Printf("main dial error=%s",err.Error())
	}
	go c1.Start()
	return c1
}

func udpClient() client.ChatClient {
	var c1 client.ChatClient
	c1 = client.NewUdpClient()
	err := c1.Dial("192.168.56.101:2222")
	if err != nil {
		log.Printf("main dial error=%s",err.Error())
	}
	go c1.Start()
	return c1
}

func tcpClient() client.ChatClient {
	var c1 client.ChatClient
	c1 = client.NewClient()
	var err = c1.Dial("192.168.56.101:1111")
	if err != nil {
		log.Printf("main dial error=%s",err.Error())
	}
	go c1.Start()
	return c1
}

func runMessages(c1 *client.ChatClient, sentMessages int, currentSum float64,
	             clientType string, clientName string, shouldRead bool,
	             times *[10000]float64) {
	const defaultValue  = 10000
	log.Printf("sending %d messages",defaultValue)
	var sum float64
	var total int
	total = sentMessages
	sum = currentSum
	var forcefulBreak = false
	var delay = 1*time.Millisecond
	for i := 0; i < defaultValue; i++ {
		var t1 = time.Now()
		err := (*c1).SendMessage(fmt.Sprintf("%d", i))
		if err != nil {
			log.Printf("error=%s",err.Error())
			if err == io.EOF {
				i = defaultValue
				forcefulBreak = true
			}
		} else {
			time.Sleep(delay)
			t1 = t1.Add(delay)
			times[i] = float64(time.Since(t1).Nanoseconds())
			if times[i] == 0 {
				println("Eitcha deu um 0")
				forcefulBreak = true
				i=defaultValue
			} else {
				total += 1
				sum += times[i]
			}
		}
	}
	if !forcefulBreak {
		println("sent")
		if shouldRead {
			println("reading")
			var t1 = time.Now()
			count := 0
			for i := range (*c1).Incoming() {
				if count % 10000 == 0 {
					println(count, i.Message)
					println("total ",total)
				}
				count+=1
				if count >= total*5 {
					break;
				}
			}
			println("sai do for")
			time.Sleep(delay)
			t1 = t1.Add(delay)
			var tReadFinal = float64(time.Since(t1).Nanoseconds())
			if tReadFinal == 0 {
				println("Eitcha deu um 0 na leitura")
			}
			sum += float64(time.Since(t1).Nanoseconds())
			println("read "+strconv.Itoa(count)+" messages")
		}
		println("read")
		calculateMeanAndSd(clientType, clientName, total, times, sum)
		println("finishing runmessages")
	} else {
		println("forceful break")
		var NewClient = createClient(clientType)
		go runMessages(&NewClient, total, sum, clientType, clientName, shouldRead, times)
	}
}

func calculateMeanAndSd(clientType string, clientName string,
	                    total int, times* [10000]float64, sum float64) {
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
	writeToFile(clientType, clientName, mean, sd, total)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}


func writeToFile(clientType string, clientName string, mean float64, sd float64, total int) {
	if !fileExists("output.txt") {
		//Write first line
		err := ioutil.WriteFile("temp.txt", []byte("clientType, clientName, mean, sd, total\n"), 0644)
		if err != nil {
			log.Fatal(err)
		}

	}
	//Append second line
	file, err := os.OpenFile("temp.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	if _, err := file.WriteString(clientType + ", " + clientName + ", " +
		FloatToString(mean)+", "+FloatToString(sd)+", "+strconv.Itoa(total)); err != nil {
		log.Fatal(err)
	}
}

func createClient(clientType string) client.ChatClient {
	if "rpc" == clientType {
		return rpcClient()
	} else if "udp" == clientType {
		return udpClient()
	} else if "tcp" == clientType {
		return tcpClient()
	} else {
		return rabbitClient(clientType)
	}
}

func runFiveClients(clientType string) {
	var shouldRead bool
	if clientType != "udp" {
		shouldRead = true
	} else {
		shouldRead = false
	}
	var c1times [10000]float64
	var c2times [10000]float64
	var c3times [10000]float64
	var c4times [10000]float64
	var c5times [10000]float64


	log.Printf("Com 5 clientes")
	log.Printf(clientType)
	var c1,c2,c3,c4,c5 client.ChatClient
	if clientType != "rabbit" {
		c1 = createClient(clientType)
		c2 = createClient(clientType)
		c3 = createClient(clientType)
		c4 = createClient(clientType)
		c5 = createClient(clientType)
	} else {
		c1 = createClient("c1")
		c2 = createClient("c2")
		c3 = createClient("c3")
		c4 = createClient("c4")
		c5 = createClient("c5")
	}
	c1.SetName("c1")
	c2.SetName("c2")
	c3.SetName("c3")
	c4.SetName("c4")
	c5.SetName("c5")
	time.Sleep(1*time.Second)
	go runMessages(&c1, 0, 0, clientType, "c1", shouldRead, &c1times)
	go runMessages(&c2, 0, 0, clientType, "c2",  shouldRead,&c2times)
	go runMessages(&c3, 0, 0, clientType, "c3",  shouldRead,&c3times)
	go runMessages(&c4, 0, 0, clientType, "c4",  shouldRead,&c4times)
	go runMessages(&c5, 0, 0, clientType, "c5",  shouldRead,&c5times)
	fmt.Scanln()
	defer c1.Close()
	defer c2.Close()
	defer c3.Close()
	defer c4.Close()
	defer c5.Close()
}

func main() {
	runFiveClients("rpc")
	fmt.Scanln()
	runFiveClients("tcp")
	fmt.Scanln()
	runFiveClients("udp")
	fmt.Scanln()
	runFiveClients("rabbit")
	fmt.Scanln()
}
