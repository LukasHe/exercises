package main

import "./Logic"
import "./NetworkModule"
import "time"
import "strconv"


func main(){


    newOrderChan := make(chan string,10)
	doneOrderChan := make(chan string,10)
	bidChan := make(chan string,10)
	sendChan := make(chan string,10)
	selfOrderChan := make(chan string,10)


	Logic.LogicInit(newOrderChan, doneOrderChan, bidChan, sendChan, selfOrderChan)
	NetworkModule.NetworkInit(sendChan, newOrderChan, doneOrderChan, bidChan)
	time.Sleep(100*time.Millisecond)
	sendChan <- "N" + "_" + strconv.Itoa(int(time.Now().UnixNano())) + "_" + "1"
	// sendChan <- "N" + strconv.Itoa(int(time.Now().UnixNano())) + "2"
	// sendChan <- "N" + strconv.Itoa(int(time.Now().UnixNano())) + "3"


	keepAlive := make(chan int)
	<- keepAlive
	
}

