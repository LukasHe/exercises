package main

import (
	"./Driver"
 	"./HardwareControll"
 	"./NetworkModule"
 	"./Logic"
 	//"time"
 	//"strconv"
 )

func main(){
	ledOnChan := make(chan int, 10)
	ledOffChan := make(chan int,10)
	sensorChan := make(chan int, 10)
	keepAlive := make(chan int)
	internalOrderChan := make(chan int,10)
	
	motorDirChan := make(chan string, 10)
	buttonChan := make(chan string, 10)
	selfOrderChan := make(chan string)
	sendChan := make(chan string, 10)
	doneOrderChan := make(chan string, 10)
	newOrderChan := make(chan string, 10)
	bidChan := make(chan string, 10)



	Logic.LogicInit(newOrderChan, doneOrderChan, bidChan, sendChan, selfOrderChan, internalOrderChan)

	Driver.DriverInit(ledOnChan, ledOffChan, sensorChan, motorDirChan, buttonChan)

	HardwareControll.HardwareControllInit(ledOnChan, ledOffChan, sensorChan , internalOrderChan, 
	motorDirChan, buttonChan, selfOrderChan, sendChan)

	NetworkModule.NetworkInit(sendChan, newOrderChan, doneOrderChan, bidChan)


	<-keepAlive
}