package main

import (
"./NetworkModule"
"fmt"
"time"
"strconv"
"os/exec"
)

func main(){
	sendChan := make(chan string, 10)
	recChan := make(chan string, 10)
	killChan := make(chan bool, 1)
	countString := "0"
	countInt := 0
	go NetworkModule.Sender(sendChan)
	go NetworkModule.Receiver(recChan, killChan)
	

	BreakHandle:
	for {
		time.Sleep(time.Millisecond * 1000)
		select {
			case countString = <-recChan:
				//fmt.Println("MainRec",countString)
				
			default:
				break BreakHandle
				
		}
				
	}

	killChan <- true
	
	cmd := exec.Command("go run main.go")
	cmd.Start()
	
	//fmt.Println("String:",countString)
	countInt, err := strconv.Atoi(countString)
	if err != nil {
        // handle error
        fmt.Println(err)
        }
	//fmt.Println("INT:",countInt)
	for {
		countInt++
		sendChan <- strconv.Itoa(countInt)
		fmt.Println("We waited for", countInt, "seconds.")
		fmt.Println("That are",countInt/60,"minutes.")
		time.Sleep(time.Millisecond * 1000)
	}
}
