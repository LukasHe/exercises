package main

import "./driver"
import "fmt"

func main(){
	ledOnChan := make(chan int, 10)
	ledOffChan := make(chan int,10)
	keepAlive := make(chan int)
	driver.DriverInit(ledOnChan, ledOffChan)
	fmt.Println(LIGHT_DOWN3)
	//ledOnChan <- LIGHT_DOWN3
	<- keepAlive
	
}

