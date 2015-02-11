package driver
/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "elev.h"
#include "channels.h"
*/
import "C"


func DriverTest(){


C.elev_init()
C.elev_set_motor_direction(DIRN_DOWN)


}