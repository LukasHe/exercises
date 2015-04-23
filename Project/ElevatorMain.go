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

	sensorChan := make(chan int, 10)
	keepAlive := make(chan int)
	internalOrderChan := make(chan int,10)
	
	ledOnChan := make(chan string, 10)
	ledOffChan := make(chan string,10)
	motorDirChan := make(chan string, 10)
	buttonChan := make(chan string, 10)
	selfOrderChan := make(chan string)
	sendChan := make(chan string, 10)
	doneOrderChan := make(chan string, 10)
	newOrderChan := make(chan string, 10)
	bidChan := make(chan string, 10)



	Logic.LogicInit(newOrderChan, doneOrderChan, bidChan, sendChan, selfOrderChan, ledOffChan, internalOrderChan)

	Driver.DriverInit(sensorChan, ledOnChan, ledOffChan, motorDirChan, buttonChan)

	HardwareControll.HardwareControllInit(sensorChan , internalOrderChan, ledOnChan, ledOffChan, motorDirChan, buttonChan, selfOrderChan, sendChan)

	NetworkModule.NetworkInit(sendChan, newOrderChan, doneOrderChan, bidChan)


	<-keepAlive
}