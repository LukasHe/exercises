package main

import "./Driver"
//import "fmt"
//import "time"

func main(){
	ledOnChan := make(chan int, 10)
	ledOffChan := make(chan int,10)
	sensorChan := make(chan int, 10)
	keepAlive := make(chan int)
	
	motorDirChan := make(chan string, 10)
	buttonChan := make(chan string, 10)
	
	driver.DriverInit(ledOnChan, ledOffChan, sensorChan, motorDirChan, buttonChan)
	
	ledOnChan <- driver.LIGHT_DOWN3
	//motorDirChan <- "DOWN"
	//time.Sleep(2000 * time.Millisecond)
	//motorDirChan <- "UP"
	//time.Sleep(2000 * time.Millisecond)
	//motorDirChan <- "STOP"
	
	<- keepAlive
	
}

