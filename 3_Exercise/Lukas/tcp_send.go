package main

import 
(
    "fmt"
    "runtime"
    "net"
)






func main(){
    runtime.GOMAXPROCS(runtime.NumCPU())
          
    serverAddr :="dieholzkatze.servebeer.com:8888"

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
		fmt.Println(err)
	} 
	fmt.Println(string(reply[:]))


    fmt.Println("Inside main! This is i: ");
}