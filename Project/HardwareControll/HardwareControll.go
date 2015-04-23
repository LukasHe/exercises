package HardwareControll

import (
	"time"
	"fmt"
	"strconv"
	"strings"
	".././NetworkModule"
)

func HardwareControllInit(sensorChan , internalOrderChan chan int, ledOnChan, ledOffChan, motorDirChan, buttonChan, selfOrderChan, sendChan chan string){
	currentFloor := 2

	//This forces the elevator to start at floor 0
	for currentFloor != 1{
		motorDirChan <- "DOWN"
		currentFloor = <- sensorChan
	}

	motorDirChan <- "STOP"

	go hardwareControll(sensorChan, internalOrderChan, ledOnChan, ledOffChan, motorDirChan, buttonChan, selfOrderChan,sendChan, currentFloor)

}

func hardwareControll(sensorChan, internalOrderChan chan int, ledOnChan, ledOffChan, motorDirChan, buttonChan, selfOrderChan, sendChan chan string, currentFloor int){
	var nextFloorOrder int
	var nextOrder string
	var timeStamp int
	var nextOrderFloorDir string
	var origIP string


	//This infinity for handles newOrder's and changes the motorDir according to where you need
	//to go. It also sends the completed order to the sendChan to be broadcasted.
	for{
		select {
			case nextOrder = <- selfOrderChan:

				// fmt.Println("New order! ", nextOrder)
				
				timeStamp, origIP, nextOrderFloorDir, _ = NetworkModule.SplitMessage(nextOrder)

				nextFloorOrder,_ = strconv.Atoi(string(nextOrderFloorDir[0]))
				// fmt.Println("current: ", currentFloor)
				// fmt.Println("Floor: ", nextFloorOrder)

				if currentFloor < nextFloorOrder && nextFloorOrder <= 4{ //Maybe change to a MAXFLOOR
					ledOnChan <- "LIGHT_DOOR_OPEN"
					time.Sleep(1000*time.Millisecond)
					ledOffChan <- "LIGHT_DOOR_OPEN"
					motorDirChan <- "UP"
				} else if currentFloor > nextFloorOrder && nextFloorOrder > 0{
					ledOnChan <- "LIGHT_DOOR_OPEN"
					time.Sleep(1000*time.Millisecond)
					ledOffChan <- "LIGHT_DOOR_OPEN"
					motorDirChan <- "DOWN"
				} else if currentFloor == nextFloorOrder {
					motorDirChan <- "STOP"
					fmt.Println("current: ", currentFloor, "nextFloorOrder: ", nextFloorOrder)
					sendChan <- "D" + "_" + strconv.Itoa(timeStamp) + "_" + origIP + "_" + nextOrderFloorDir + "_" + "*"
				} else {
					motorDirChan <- "STOP"
				}

			case pressedButton := <- buttonChan:
				splitButtons := strings.Split(pressedButton, "_")
				switch splitButtons[1]{
				case "COMMAND":
					sendChan <- "N" + "_" + strconv.Itoa(int(time.Now().UnixNano())) + "_" + NetworkModule.GetOwnIP() + "_" + splitButtons[2] + "I"+ "_" + "I"
					ledOnChan <- "LIGHT_" + splitButtons[1] + splitButtons[2]
				case "UP":
					sendChan <- "N" + "_" + strconv.Itoa(int(time.Now().UnixNano())) + "_" + NetworkModule.GetOwnIP() + "_" + splitButtons[2] + "U" + "_" + NetworkModule.GetOwnIP()
					ledOnChan <- "LIGHT_" + splitButtons[1] + splitButtons[2]
				case "DOWN":
					sendChan <- "N" + "_" + strconv.Itoa(int(time.Now().UnixNano())) + "_" + NetworkModule.GetOwnIP() + "_" + splitButtons[2] + "D" + "_" + NetworkModule.GetOwnIP()
					ledOnChan <- "LIGHT_" + splitButtons[1] + splitButtons[2]
				case "STOP":
					motorDirChan <- "STOP"
					ledOnChan <- "LIGHT_STOP"
					fmt.Println("Stop Button was pressed!!!!!!!!!")
				}


			case currentFloor = <- sensorChan:

				if currentFloor == nextFloorOrder{
					motorDirChan <- "STOP"
					fmt.Println("current: ", currentFloor, "nextFloorOrder: " ,nextFloorOrder)
					sendChan <- "D" + "_" + strconv.Itoa(timeStamp) + "_" + origIP + "_" + nextOrderFloorDir + "_" + "*"
				}
			default:
				time.Sleep(10*time.Millisecond)

		}
	}
}