package driver
/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"
import "fmt"
import "time"

const PORT4             =  3
const OBSTRUCTION       =  (0x300+23)
const STOP              =  (0x300+22)
const BUTTON_COMMAND1   =  (0x300+21)
const BUTTON_COMMAND2   =  (0x300+20)
const BUTTON_COMMAND3   =  (0x300+19)
const BUTTON_COMMAND4   =  (0x300+18)
const BUTTON_UP1        =  (0x300+17)
const BUTTON_UP2        =  (0x300+16)

const PORT1             =  2
const BUTTON_DOWN2      =  (0x200+0)
const BUTTON_UP3        =  (0x200+1)
const BUTTON_DOWN3      =  (0x200+2)
const BUTTON_DOWN4      =  (0x200+3)
const SENSOR_FLOOR1     =  (0x200+4)
const SENSOR_FLOOR2     =  (0x200+5)
const SENSOR_FLOOR3     =  (0x200+6)
const SENSOR_FLOOR4     =  (0x200+7)

const PORT3             =  3
const MOTORDIR          =  (0x300+15)
const LIGHT_STOP        =  (0x300+14)
const LIGHT_COMMAND1    =  (0x300+13)
const LIGHT_COMMAND2    =  (0x300+12)
const LIGHT_COMMAND3    =  (0x300+11)
const LIGHT_COMMAND4    =  (0x300+10)
const LIGHT_UP1         =  (0x300+9)
const LIGHT_UP2         =  (0x300+8)

const PORT2             =  3
const LIGHT_DOWN2       =  (0x300+7)
const LIGHT_UP3         =  (0x300+6)
const LIGHT_DOWN3       =  (0x300+5)
const LIGHT_DOWN4       =  (0x300+4)
const LIGHT_DOOR_OPEN   =  (0x300+3)
const LIGHT_FLOOR_IND2  =  (0x300+1)
const LIGHT_FLOOR_IND1  =  (0x300+0)


const PORT0             =  1
const MOTOR             =  (0x100+0)

const BUTTON_DOWN1      =  -1
const BUTTON_UP4        =  -1
const LIGHT_DOWN1       =  -1
const LIGHT_UP4         =  -1

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
		}
	}
}
			
