package main

import 
(
    "fmt"
    "runtime"
    "net"
)

func main(){
    runtime.GOMAXPROCS(runtime.NumCPU())
          
    serverAddr :="129.241.187.136:33546"

    tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
    if err != nil {
    fmt.Println("error resolve TCP-address")
    fmt.Println(err)
    }

    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
		fmt.Println("Dial failed:")
		fmt.Println(err)
	}

 	reply := make([]byte, 1024)
 
	_, err = conn.Read(reply)
	if err != nil {
		println("Write to server failed:")
		println(err)
	} 
	fmt.Println(string(reply[:]))

	sendString := make([]byte, 1024)
	sendString = []byte("Connect to: 129.241.187.151:2222\x00")
	_, err = conn.Write(sendString)
	if err != nil {
		println("Write to server failed:")
		println(err)
	} 

	n, err := conn.Read(reply)
	if err != nil {
		println("Write to server failed:")
		println(err)
	} 
	fmt.Println(string(reply[:]),n)


}