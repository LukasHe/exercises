package NetworkModule

import (
"fmt"
"net"
"time"
)

func sender(sendChannel chan string) {
	message := ""
	broadcastAddr := "129.241.187.255:22022"

	udpBroadcast, err := net.Dial("udp", broadcastAddr)
	if err != nil {
		fmt.Println("Could not resolve udp address or connect to it.")
		fmt.Println(err)
		return
	}
	
	defer udpBroadcast.Close()

	for {
		time.Sleep(100*time.Millisecond)
		message =<- sendChannel
		_, err = udpBroadcast.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing data to server:", broadcastAddr)
			fmt.Println(err)
			return
		}
		fmt.Println(message)
	}		
}

func Receiver() {
	
	var receiveBuf [1024]byte

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

	for {
		rlen, remote, err := udpReceive.ReadFromUDP(receiveBuf[:])
		if err != nil {
			fmt.Println("error reading from UDP", remote)
			fmt.Println(err)
			return
		}
		fmt.Println("Recived", rlen ,"Byte from", remote, ".")
		fmt.Println("The message is:",string(receiveBuf[:]))
	}

}