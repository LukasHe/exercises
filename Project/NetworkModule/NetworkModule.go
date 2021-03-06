package NetworkModule

import (
"fmt"
"net"
"time"
"strings"
"strconv"
)

//Initialize function that spawn two threads.
func NetworkInit(sendChan, newOrderChan, doneOrderChan, bidChan chan string){

	go sender(sendChan)
	go receiver(newOrderChan, doneOrderChan, bidChan)
}

//Establishes a UDP connection by using .Dial and broadcasts all 
//incoming messages to broadcastAddr.
func sender(sendChan chan string) {

	broadcastAddr := "129.241.187.255:22021"
	message := ""
	udpBroadcast, err := net.Dial("udp", broadcastAddr)

	if err != nil {
		fmt.Println("Could not resolve udp address or connect to it.")
		fmt.Println(err)
		return
	}
	defer udpBroadcast.Close()

	for {
		message =<- sendChan
		_, err = udpBroadcast.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing data to server:", broadcastAddr)
			fmt.Println(err)
			return
		}
		fmt.Println("Send:",message)
		time.Sleep(10*time.Millisecond)
	}		
}


//Listen for UDP messages on given port and distribute them according
//to the first element in the message. 
func receiver(newOrderChan, doneOrderChan, bidChan chan string) {

	var receiveBuf [1024]byte
	udpAddr, err := net.ResolveUDPAddr("udp", ":22021")

	if err != nil {
		fmt.Println("Error resolve UDP-address")
		fmt.Println(err)
		return
	}

	udpReceive, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		fmt.Println("Error listen to UDP")
		fmt.Println(err)
		return
	}

	defer udpReceive.Close()

	for {
		rlen, remoteIP, err := udpReceive.ReadFromUDP(receiveBuf[:])

		if err != nil {
			fmt.Println("error reading from UDP", remoteIP)
			fmt.Println(err)
			return
		}

		switch {
			case "D" == string(receiveBuf[0]):
				doneOrderChan <- string(receiveBuf[2:rlen]) 
			case "N" == string(receiveBuf[0]):
				newOrderChan<- string(receiveBuf[2:rlen])
			case "B" == string(receiveBuf[0]):
				bidChan <- string(receiveBuf[2:rlen])
			default:
				time.Sleep(10*time.Millisecond)	
		}
	}
}

//Global function that returns the IP of the caller
func GetOwnIP() (string) {

	allAddrs, _ := net.InterfaceAddrs()
	v4Addr := strings.Split(allAddrs[1].String(), "/")
	completeIP := strings.Split(v4Addr[0],".")
	return completeIP[3]
}

//Global function that split up input string and returns the 
//seperate elements. 
func SplitMessage(message string) (int, string, string, string) {

	splitMsg := strings.Split(message, "_")
	time, originIP, data, remoteIP := splitMsg[0], splitMsg[1], splitMsg[2], splitMsg[3]
	timeStamp, _ := strconv.Atoi(time)
	return timeStamp, originIP, data, remoteIP
}