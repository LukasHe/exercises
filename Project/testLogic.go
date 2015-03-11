package main

import "./Logic"
import "./NetworkModule"
import "time"


func main(){


    newOrderChan := make(chan string,10)
	doneOrderChan := make(chan string,10)
	bidChan := make(chan string,10)
	sendChan := make(chan string,10)

	Logic.LogicInit(newOrderChan, doneOrderChan, bidChan)
	NetworkModule.NetworkInit(sendChan, newOrderChan, doneOrderChan, bidChan)
	time.Sleep(100*time.Millisecond)
	sendChan <- "DdoneOrder"
	sendChan <- "NnewOrder"
	sendChan <- "Bbid"
	sendChan <- "DdoneOrder"

	keepAlive := make(chan int)
	<- keepAlive
	
}

