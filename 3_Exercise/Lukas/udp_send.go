package main

import (
"fmt"
"net"
"time"
"strconv"
)

func main() {
    fmt.Println("start main")
    // arrayToSend := []byte("This is the ? Message\n")
    var msgCounter int = 0
    serverAddr :="129.241.187.255:30018"
    fmt.Println("start Dail")
    conn, err := net.Dial("udp", serverAddr)

    if err != nil {
            fmt.Println("Could not resolve udp address or connect to it.")
            fmt.Println(err)
            return
    }

    fmt.Println("Connected to server at ", serverAddr)

    defer conn.Close()

    fmt.Println("About to write to connection")

    for {
        time.Sleep(1000*time.Millisecond)
        
        n, err := conn.Write([]byte("Message "+ strconv.Itoa(msgCounter)))
        if err != nil {
            fmt.Println("error writing data to server", serverAddr)
            fmt.Println(err)
            return
        }
        msgCounter++

        if n > 0 {
            fmt.Println("Wrote ",n, " bytes to server at ", serverAddr)
        }
    }

}