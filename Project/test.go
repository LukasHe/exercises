package main

import (
	"./Driver"
 	"./HardwareControll"
 	"./NetworkModule"
 	"time"
 	"strconv"
 )

func main(){
	ledOnChan := make(chan int, 10)
	ledOffChan := make(chan int,10)
	sensorChan := make(chan int, 10)
	keepAlive := make(chan int)
	
	motorDirChan := make(chan string, 10)
	buttonChan := make(chan string, 10)
	selfOrderChan := make(chan string)

	sendChan := make(chan string, 10)
	doneOrderChan := make(chan string, 10)
	newOrderChan := make(chan string, 10)
	bidChan := make(chan string, 10)

	Driver.DriverInit(ledOnChan, ledOffChan, sensorChan, motorDirChan, buttonChan)
	HardwareControll.HardwareControllInit(ledOnChan, ledOffChan, sensorChan, motorDirChan, 
		buttonChan, selfOrderChan, sendChan)
	NetworkModule.NetworkInit(sendChan, newOrderChan, doneOrderChan, bidChan)


	time.Sleep(5000*time.Millisecond)
	selfOrderChan <- "D" + strconv.Itoa(int(time.Now().UnixNano())) + "3"
	time.Sleep(10000*time.Millisecond)
	selfOrderChan <- "D" + strconv.Itoa(int(time.Now().UnixNano())) + "2"
		time.Sleep(10000*time.Millisecond)
	selfOrderChan <- "D" + strconv.Itoa(int(time.Now().UnixNano())) + "2"
	time.Sleep(10000*time.Millisecond)
	selfOrderChan <- "D" + strconv.Itoa(int(time.Now().UnixNano())) + "0"
	time.Sleep(10000*time.Millisecond)
	selfOrderChan <- "D" + strconv.Itoa(int(time.Now().UnixNano())) + "1"
	time.Sleep(10000*time.Millisecond)
	selfOrderChan <- "D" + strconv.Itoa(int(time.Now().UnixNano())) + "0"


	//ledOnChan <- Driver.LIGHT_DOWN3


	//motorDirChan <- "DOWN"
	//time.Sleep(2000 * time.Millisecond)
	//motorDirChan <- "UP"
	//time.Sleep(2000 * time.Millisecond)
	//motorDirChan <- "STOP"
	
	<- keepAlive
	
}

