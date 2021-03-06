package Driver
/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"
import "fmt"
import "time"

//Initialize the Driver.
func DriverInit(sensorChan chan int, ledOnChan, ledOffChan, motorDirChan, buttonChan chan string){
		
	lightArray :=[...] int {
		LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1,
		LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2,
		LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3,
		LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4,
		LIGHT_STOP, LIGHT_DOOR_OPEN,
		LIGHT_FLOOR_IND1, LIGHT_FLOOR_IND2}
	
	// Init hardware
	if C.io_init() != 1 {
		//make error
		fmt.Println("error in io_init")
	}

	// Set everything to zero
	for i := 0; i < len(lightArray); i++ {
		C.io_clear_bit(C.int(lightArray[i]))
	}
	time.Sleep(100*time.Millisecond)
	
	// Start the Hardwaredriver as thread
	go driver(sensorChan, motorDirChan, ledOnChan, ledOffChan, buttonChan)
}


func driver(sensorChan chan int, motorDirChan, ledOnChan, ledOffChan, buttonChan chan string){
	
	sensorArray :=[]int{
	SENSOR_FLOOR1, SENSOR_FLOOR2, 
	SENSOR_FLOOR3, SENSOR_FLOOR4}

	oldSensorValue := []int{0,0,0,0}
	currentSensorValue := 0

	buttonNameArray :=[]string{
		"BUTTON_COMMAND_1", "BUTTON_COMMAND_2", 
		"BUTTON_COMMAND_3", "BUTTON_COMMAND_4",
		"BUTTON_UP_1", "BUTTON_UP_2",
		"BUTTON_UP_3", "BUTTON_DOWN_2",
		"BUTTON_DOWN_3", "BUTTON_DOWN_4",
		"BUTTON_STOP_0"}

	buttonArray :=[]int{
		BUTTON_COMMAND1, BUTTON_COMMAND2, 
		BUTTON_COMMAND3, BUTTON_COMMAND4,
		BUTTON_UP1, BUTTON_UP2,
		BUTTON_UP3, BUTTON_DOWN2,
		BUTTON_DOWN3, BUTTON_DOWN4,
		BUTTON_STOP} 

	oldButtonValue := []int{0,0,0,0,0,0,0,0,0,0,0}
	currentButtonValue := 0


	buttonLightMap := map[string] int{
		"LIGHT_STOP":		LIGHT_STOP,
		"LIGHT_COMMAND1":  	LIGHT_COMMAND1,
		"LIGHT_COMMAND2":	LIGHT_COMMAND2,
		"LIGHT_COMMAND3":   LIGHT_COMMAND3,
		"LIGHT_COMMAND4":   LIGHT_COMMAND4,
		"LIGHT_UP1":        LIGHT_UP1,
		"LIGHT_UP2":        LIGHT_UP2,
		"LIGHT_DOWN2":      LIGHT_DOWN2,
		"LIGHT_UP3":        LIGHT_UP3,
		"LIGHT_DOWN3":      LIGHT_DOWN3,
		"LIGHT_DOWN4":      LIGHT_DOWN4,
		"LIGHT_DOOR_OPEN":  LIGHT_DOOR_OPEN,
	}

	for{
	
		//Check LED- and Motor-Channel for updates and apply them to the elevator.
		select {

			case ledOnOrder := <-ledOnChan:
					C.io_set_bit(C.int(buttonLightMap[ledOnOrder]))

			case ledOffOrder := <- ledOffChan:
					C.io_clear_bit(C.int(buttonLightMap[ledOffOrder]))

			case motorDir := <-motorDirChan:

					if motorDir == "UP"{
						C.io_clear_bit(MOTORDIR);
						C.io_write_analog(MOTOR, 2800);

					} else if motorDir == "DOWN"{
						C.io_set_bit(MOTORDIR);
						C.io_write_analog(MOTOR, 2800);

					} else{
						C.io_write_analog(MOTOR, 0);
					}

			default:
				time.Sleep(10*time.Millisecond)
		}
		
		broadcastSensorValue(sensorChan, sensorArray, oldSensorValue, currentSensorValue)
		broadcastButtonPress(buttonChan, buttonNameArray, buttonArray, oldButtonValue, currentButtonValue)

	}
}
	
//Set Floor LEDs.		
func setFloorLED(floor int){

	switch floor{

		case 1:
			C.io_clear_bit(LIGHT_FLOOR_IND1)
			C.io_clear_bit(LIGHT_FLOOR_IND2)

		case 2:
			C.io_clear_bit(LIGHT_FLOOR_IND1)
			C.io_set_bit(LIGHT_FLOOR_IND2)

		case 3:
			C.io_set_bit(LIGHT_FLOOR_IND1)
			C.io_clear_bit(LIGHT_FLOOR_IND2)

		case 4:
			C.io_set_bit(LIGHT_FLOOR_IND1)
			C.io_set_bit(LIGHT_FLOOR_IND2)
	}
}

//Go through all Sensors and send index of the active ones on channel.
func broadcastSensorValue(sensorChan chan int, sensorArray, oldSensorValue []int, currentSensorValue int){

	for index, floorSensor := range sensorArray{
	currentSensorValue = int(C.io_read_bit(C.int(floorSensor)))

		if oldSensorValue[index] != currentSensorValue{
			oldSensorValue[index] = currentSensorValue

			if currentSensorValue == 1 {
				floor := index + 1
				sensorChan <- floor
				setFloorLED(floor)
			}
		}
	}
}

//Go throug all Buttons and send buttonspress on channel.
func broadcastButtonPress(buttonChan chan string, buttonNameArray []string , buttonArray, oldButtonValue []int, currentButtonValue int){

	for index, button := range buttonArray{
		currentButtonValue = int(C.io_read_bit(C.int(button)))

		if oldButtonValue[index] != currentButtonValue{
			oldButtonValue[index] = currentButtonValue

			if currentButtonValue == 1 {
				buttonChan <- buttonNameArray[index]
			}
		}
	}	
}