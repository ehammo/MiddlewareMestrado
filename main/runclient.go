package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	d "../distribution"
	n "../naming"
)

func startClient(id int) *d.ClientProxy {
	namingProxy := n.NewNamingProxy("localhost:1243")

	aor := namingProxy.LookUp("Vanet")
	fmt.Println("Received aor:")
	fmt.Println(aor)
	var c = d.NewClientProxy(aor, id)
	go c.Start()
	return c
}

func threeBreakingCars() {
	log.Printf("twoBreakingCars")
	var c1, c2, c3, c4, c5 *d.ClientProxy
	c1 = setupAClient("0")
	c2 = setupAClient("1")
	c3 = setupAClient("2")
	c4 = setupAClient("3")
	c5 = setupAClient("4")
	c1.RegisterOnLane("lane1")
	c2.RegisterOnLane("lane1")
	c3.RegisterOnLane("lane2")
	c4.RegisterOnLane("lane2")
	c5.RegisterOnLane("lane2")
	time.Sleep(5 * time.Second)
	fmt.Println("5 seconds to go")
	time.Sleep(5 * time.Second)
	c1.BroadcastEvent("lane2")
}
func setupAClient(id string) *d.ClientProxy {
	intid, _ := strconv.Atoi(id)
	c1 := startClient(intid)
	c1.RegisterKey()
	go c1.Start()
	return c1
}

func trimClientIdPlusC(clientId string) string {
	toint, _ := strconv.Atoi(clientId)
	clientId = "c" + strconv.Itoa(toint)
	return clientId
}

func writeToFile(clientId string, data string) {
	clientId = trimClientIdPlusC(clientId)
	var filename = clientId + "-result.txt"

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		err := ioutil.WriteFile(filename, []byte(data), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer file.Close()
	if _, err := file.WriteString(data); err != nil {
		log.Fatal(err)
	}

}

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func main() {
	var total, sd float64 = 0, 0
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter uniqueId for this Car: ")
	id, _ := reader.ReadString('\n')
	t1 := time.Now()
	c1 := setupAClient(id)
	t2 := time.Since(t1).Nanoseconds()
	t2InMili := float64(t2) / 1000000
	setupTime := "Time to setup: " + FloatToString(t2InMili) + " mili"
	writeToFile(id, setupTime+"\n")
	fmt.Println(setupTime)
	fmt.Println("Stress use case. Send 1000 messages to server")

	var times []float64 = make([]float64, 1000, 1000)
	for i := 0; i < 1000; i++ {
		midTime := time.Now()
		if i == 0 {
			c1.RegisterOnLane("lane1")
		} else if i%2 == 0 {
			c1.ChangeLane("lane2")
		} else {
			c1.ChangeLane("lane1")
		}
		midTime2 := time.Since(midTime).Nanoseconds()
		times[i] = float64(midTime2)
		total += times[i]
	}
	mean := total / 1000
	for i := 0; i < 1000; i++ {
		if times[i] > 0 {
			sd += math.Pow(times[i]-mean, 2)
		}
	}
	sd = math.Sqrt(sd / float64(total))
	t2InMili = float64(t2) / 1000000
	fmt.Println("Total time ", t2InMili, " milisecond")
	meanInMili := mean / 1000000
	sdInMili := sd / 1000000
	fmt.Println("Mean ", meanInMili, " mili")
	fmt.Println("Sd ", sdInMili, " mili")
	meanSd := "Mean: " + FloatToString(meanInMili) + " mili\nSd: " + FloatToString(sdInMili) + " mili"
	writeToFile(id, meanSd+"\n")
	fmt.Scanln()
	// twoBreakingCars()
	// fmt.Scanln()
}
