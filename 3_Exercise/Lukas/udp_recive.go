// Go 1.2
// go run udp_recive.go

package main

import 
(
    "fmt"
    "runtime"
    "net"
)

var messages = make(chan bool)

func udp_receive(){

    fmt.Println("receive start")
    var buf [1024]byte

    addr, err := net.ResolveUDPAddr("udp", ":20018")
    if err != nil {
        fmt.Println("error resolve UDP-address")
        fmt.Println(err)
        return
    }

    sock, err := net.ListenUDP("udp", addr)
    if err != nil {
        fmt.Println("error listen to UDP")
        fmt.Println(err)
        return
    }

    defer sock.Close()

    for {
        rlen, remote, err := sock.ReadFromUDP(buf[:])
        if err != nil {
            fmt.Println("error reading from UDP", remote)
            fmt.Println(err)
            return
        }
        fmt.Println("Recived", rlen ,"Byte from", remote, ".")
        fmt.Println("The message is:",string(buf[:]))
    }
    fmt.Println("receive done")
}


func main(){
    runtime.GOMAXPROCS(runtime.NumCPU())
   
    udp_receive()         
    

    fmt.Println("Inside main! This is i: ");
}