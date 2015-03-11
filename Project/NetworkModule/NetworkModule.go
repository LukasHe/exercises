package NetworkModule

import (
"fmt"
"net"
"time"
)

func NetworkInit(sendChannel, newOrderChan, doneOrderChan, bidChan chan string){
	go sender(sendChannel)
	go receiver(newOrderChan, doneOrderChan, bidChan)
}

func sender(sendChannel chan string) {
	broadcastAddr := "129.241.187.255:22022"
	message := ""
	
	//Connect to broadcastAddr using UDP.
	udpBroadcast, err := net.Dial("udp", broadcastAddr)

	if err != nil {
		fmt.Println("Could not resolve udp address or connect to it.")
		fmt.Println(err)
		return
	}
	defer udpBroadcast.Close()

	for {
		message =<- sendChannel
		_, err = udpBroadcast.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing data to server:", broadcastAddr)
			fmt.Println(err)
			return
		}
		//fmt.Println("Send:",message)
		time.Sleep(10*time.Millisecond)
	}		
}

func receiver(newOrderChan, doneOrderChan, bidChan chan string) {
	var receiveBuf [1024]byte

	//Create an  UDPAddr struct and listens to the given port.
	udpAddr, err := net.ResolveUDPAddr("udp", ":22022")
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
		rlen, remote, err := udpReceive.ReadFromUDP(receiveBuf[:])
		if err != nil {
			fmt.Println("error reading from UDP", remote)
			fmt.Println(err)
			return
		}
		//fmt.Println("Recived", rlen ,"Byte from", remote, ".")
		//fmt.Println("The message is:",string(receiveBuf[:]))
		switch {
			case "D" == string(receiveBuf[0]):
				doneOrderChan <- string(receiveBuf[1:rlen])
			case "N" == string(receiveBuf[0]):
				newOrderChan<- string(receiveBuf[1:rlen])
			case "B" == string(receiveBuf[0]):
				bidChan <- string(receiveBuf[1:rlen])
			default:
				time.Sleep(10*time.Millisecond)	
		}
	}

}