package main

import 
(
    "fmt"
    "runtime"
    "net"
)

func main(){
    runtime.GOMAXPROCS(runtime.NumCPU())
          

    ln, err := net.Listen("tcp", ":2222")
    if err != nil {
    fmt.Println("error setting up listen")
    fmt.Println(err)
    }

    conn, err := ln.Accept()
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
    sendString = []byte("Patrik and Lukas msg\x00")
    _, err = conn.Write(sendString)
    if err != nil {
        println("Write to server failed:")
        println(err)
    } 

    _, err = conn.Read(reply)
    if err != nil {
        println("Write to server failed:")
        println(err)
    }
    fmt.Println(string(reply[:]))

}