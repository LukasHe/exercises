package main

import "./driver"
//import "fmt"
//import "time"

func main(){
	ledOnChan := make(chan int, 10)
	ledOffChan := make(chan int,10)
	
	motorDirChan := make(chan string, 10)
	sensorChan := make(chan int, 10)
	
	keepAlive := make(chan int)
	
	
	driver.DriverInit(ledOnChan, ledOffChan, sensorChan, motorDirChan )
	
	ledOnChan <- driver.LIGHT_DOWN3
	//motorDirChan <- "DOWN"
	//time.Sleep(2000 * time.Millisecond)
	//motorDirChan <- "UP"
	//time.Sleep(2000 * time.Millisecond)
	//motorDirChan <- "STOP"
	
	<- keepAlive
	
}

