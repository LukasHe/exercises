package HardwareControll

import "time"
import "fmt"
import "strconv"
//import "reflect"

func HardwareControllInit(ledOnChan, ledOffChan, sensorChan chan int, 
	motorDirChan, buttonChan, selfOrderChan, sendChan chan string){
	currentFloor := 1

	//This forces the elevator to start at floor 0
	for currentFloor != 0{
		motorDirChan <- "DOWN"
		currentFloor = <- sensorChan
	}
	motorDirChan <- "STOP"

	go hardwareControll(ledOnChan, ledOffChan, sensorChan, motorDirChan, buttonChan, selfOrderChan,
	 sendChan, currentFloor)

}

func hardwareControll(ledOnChan, ledOffChan, sensorChan chan int, motorDirChan, buttonChan, 
	selfOrderChan, sendChan chan string, currentFloor int){
	var nextFloorOrder int
	var nextOrder string
	//This infinity for handles newOrder's and changes the motorDir according to where you need
	//to go. It also sends the completed order to the sendChan to be broadcasted.
	for{
		select {
			case nextOrder = <- selfOrderChan:
				nextFloorOrder , _ = strconv.Atoi(string(nextOrder[20]))
				fmt.Println("Next Floor:", nextFloorOrder)
				if currentFloor < nextFloorOrder && nextFloorOrder < 4{ //Maybe change to a MAXFLOOR
					motorDirChan <- "UP"
				} else if currentFloor > nextFloorOrder && nextFloorOrder >= 0{
					motorDirChan <- "DOWN"
				} else if currentFloor == nextFloorOrder {
					motorDirChan <- "STOP"
					sendChan <- nextOrder
					time.Sleep(3000*time.Millisecond)
				} else {
					motorDirChan <- "STOP"
				}

			case currentFloor = <- sensorChan:
				if currentFloor == nextFloorOrder{
					motorDirChan <- "STOP"
					sendChan <- nextOrder
					time.Sleep(3000*time.Millisecond)
				}
			default:
				time.Sleep(10*time.Millisecond)
		}
	}
}