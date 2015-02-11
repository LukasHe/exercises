package driver
/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "channels.h"
*/
import "C"
import "fmt"

func DriverInit(){
//func DriverInit(ledChan, motorDirChan, buttonChan, sensorChan chan string){
	lightArray :=[...]C.int{
		C.LIGHT_UP1, C.LIGHT_DOWN1, C.LIGHT_COMMAND1,
    	C.LIGHT_UP2, C.LIGHT_DOWN2, C.LIGHT_COMMAND2,
    	C.LIGHT_UP3, C.LIGHT_DOWN3, C.LIGHT_COMMAND3,
    	C.LIGHT_UP4, C.LIGHT_DOWN4, C.LIGHT_COMMAND4,
    	C.LIGHT_STOP, C.LIGHT_DOOR_OPEN}
	// Init hardware
    if C.io_init() != 1 {
    	//make error
    	fmt.Println("error in io_init")
    }

    for i := 0; i < len(lightArray); i++ {
    	C.io_clear_bit(lightArray[i])
    }
    // set everything to zero

	//go driver(ledChan, motorDirChan, buttonChan, sensorChan)
}

func driver(){








}