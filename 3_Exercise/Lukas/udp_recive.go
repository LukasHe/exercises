// Go 1.2
// go run threads_in_go.go

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

    addr, err := net.ResolveUDPAddr("udp", ":2222")
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

    for {
        fmt.Println("for")
        rlen, remote, err := sock.ReadFromUDP(buf[:])
        if err != nil {
            fmt.Println("error reading from UDP", remote)
            fmt.Println(err)
            return
        }
        fmt.Println(rlen)
        fmt.Println(string(buf[:rlen]))
        if rlen > 0 {
        //messages <- true
        }
    }
    fmt.Println("receive done")
}


func main(){
    runtime.GOMAXPROCS(runtime.NumCPU())
   
    go udp_receive()         
    <- messages

    fmt.Println("Inside main! This is i: ");
}