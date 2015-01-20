package main

import (
"fmt"
"net"
"time"
)

func main() {
    arrayToSend := []byte("This is the ? Message\n")
    var msgCounter int = 0
    conn, err := net.Dial("udp", ":2222")

    if err != nil {
            fmt.Println("Could not resolve udp address or connect to it.")
            fmt.Println(err)
            return
    }

    fmt.Println("Connected to server at ", ":2222")

    defer conn.Close()

    fmt.Println("About to write to connection")

    for {
        time.Sleep(1000*time.Millisecond)
        arrayToSend[12] = byte(char(msgCounter))
        n, err := conn.Write(arrayToSend)
        if err != nil {
            fmt.Println("error writing data to server", ":2222")
            fmt.Println(err)
            return
        }
        msgCounter++

        if n > 0 {
            fmt.Println("Wrote ",n, " bytes to server at ", ":2222")
        }
    }

}