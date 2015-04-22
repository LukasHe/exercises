package HardwareControll

import (
	"time"
	"fmt"
	"strconv"
	"strings"
)
//import "reflect"

func HardwareControllInit(ledOnChan, ledOffChan, sensorChan , internalOrderChan chan int, 
	motorDirChan, buttonChan, selfOrderChan, sendChan chan string){
	currentFloor := 1

	//This forces the elevator to start at floor 0
	for currentFloor != 0{
		motorDirChan <- "DOWN"
		currentFloor = <- sensorChan
	}
	motorDirChan <- "STOP"

	go hardwareControll(ledOnChan, ledOffChan, sensorChan, internalOrderChan, motorDirChan, buttonChan, selfOrderChan,
	 sendChan, currentFloor)

}

func hardwareControll(ledOnChan, ledOffChan, sensorChan, internalOrderChan chan int, motorDirChan, buttonChan, 
	selfOrderChan, sendChan chan string, currentFloor int){
	var nextFloorOrder int
	var nextOrder string
	var nextOrderTimestamp string
	var nextOrderFloorDir string


	//This infinity for handles newOrder's and changes the motorDir according to where you need
	//to go. It also sends the completed order to the sendChan to be broadcasted.
	for{
		select {
			case nextOrder = <- selfOrderChan:
				
				nextOrderSplit := strings.Split(nextOrder, "_")
				nextOrderFloorDir,nextOrderTimestamp = nextOrderSplit[0],nextOrderSplit[1]
				nextFloorOrder := int(nextOrderFloorDir[0])

				if currentFloor < nextFloorOrder && nextFloorOrder < 4{ //Maybe change to a MAXFLOOR
					motorDirChan <- "UP"
				} else if currentFloor > nextFloorOrder && nextFloorOrder >= 0{
					motorDirChan <- "DOWN"
				} else if currentFloor == nextFloorOrder {
					motorDirChan <- "STOP"
					sendChan <- "D" + "_" + nextOrderTimestamp + "_" + nextOrderFloorDir
					time.Sleep(3000*time.Millisecond)
				} else {
					motorDirChan <- "STOP"
				}

			case pressedButton := <- buttonChan:
				splitButtons := strings.Split(pressedButton, "_")
				switch splitButtons[1]{
				case "COMMAND":
					floor,_ := strconv.Atoi(splitButtons[2])
					internalOrderChan <- floor
				case "UP":
					sendChan <- "N" + "_" + strconv.Itoa(int(time.Now().UnixNano())) + "_" + splitButtons[2] + "U" 
				case "DOWN":
					sendChan <- "N" + "_" + strconv.Itoa(int(time.Now().UnixNano())) + "_" + splitButtons[2] + "D" 
				case "STOP":
					motorDirChan <- "STOP"
					fmt.Println("Stop Button was pressed!!!!!!!!!")
				}


			case currentFloor = <- sensorChan:
				if currentFloor == nextFloorOrder{
					motorDirChan <- "STOP"
					sendChan <- "D" + "_" + nextOrderTimestamp + "_" + nextOrderFloorDir
					time.Sleep(3000*time.Millisecond)
				}
			default:
				time.Sleep(10*time.Millisecond)
		}
	}
}