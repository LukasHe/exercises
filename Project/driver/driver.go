package driver
/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"
import "fmt"
import "time"


func DriverInit(ledOnChan, ledOffChan chan int){
	//func DriverInit(ledChan, motorDirChan, buttonChan, sensorChan chan string){
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

    for i := 0; i < len(lightArray); i++ {
    	C.io_clear_bit(C.int(lightArray[i]))
    }
    time.Sleep(1000*time.Millisecond)
    // set everything to zero
fmt.Println(LIGHT_DOWN3)

	//go driver(ledChan, motorDirChan, buttonChan, sensorChan)
	go driver(ledOnChan, ledOffChan)
}

//func driver(ledChan chan string, motorDirChan chan string, buttonChan chan string, sensorChan chan string)
func driver(ledOnChan, ledOffChan chan int){
	for{
		select {
			case ledOnOrder := <-ledOnChan:
					C.io_set_bit(C.int(ledOnOrder))
			case ledOffOrder := <- ledOffChan:
					C.io_set_bit(C.int(ledOffOrder))
			case motorDir := <-motorDirChan:
					
		}
	}
}
			
