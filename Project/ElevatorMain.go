package main

import (
	"./Driver"
 	"./HardwareControll"
 	"./NetworkModule"
 	"./Logic"
 )

//The main function calles all the modules of the elevator and spawn the channels that will act as
//the main communication between the modules. This keep on running since "keepAlive" will wait.
func main(){
	sensorChan 			:= make(chan int, 10)
	internalOrderChan	:= make(chan int, 10)
	keepAlive		    := make(chan int)

	ledOnChan			:= make(chan string, 10)
	ledOffChan			:= make(chan string, 10)
	motorDirChan		:= make(chan string, 10)
	buttonChan			:= make(chan string, 10)
	selfOrderChan		:= make(chan string, 10)
	sendChan			:= make(chan string, 10)
	doneOrderChan		:= make(chan string, 10)
	newOrderChan		:= make(chan string, 10)
	bidChan				:= make(chan string, 10)

	Logic.LogicInit(newOrderChan, doneOrderChan, bidChan, sendChan, selfOrderChan, ledOffChan, internalOrderChan)
	Driver.DriverInit(sensorChan, ledOnChan, ledOffChan, motorDirChan, buttonChan)
	HardwareControll.HardwareControllInit(sensorChan , internalOrderChan, ledOnChan, ledOffChan, motorDirChan, buttonChan, selfOrderChan, sendChan)
	NetworkModule.NetworkInit(sendChan, newOrderChan, doneOrderChan, bidChan)

	<-keepAlive
}