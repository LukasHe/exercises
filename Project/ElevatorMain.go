package main
	

import (
"fmt"
"time"
"./NetworkModule"
"strconv"
)

func main(){
	sendChannel := make(chan string, 100)
	keepAlive := make(chan int)

	go NetworkModule.Sender(sendChannel)
	go NetworkModule.Receiver()
	sendChannel <- "something"
	sendChannel <- "something else"
	fmt.Println(time.Now().UnixNano())
	sendChannel <- "D" + strconv.Itoa(int(time.Now().UnixNano())) + "data"


	<-keepAlive
}