package driver

/*
#cgo CFLAGS: -std=c99
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "channels.h"
#include "elev.h"
*/
import "C"
import "source"


func Init() int {
	return int(C.elev_init())
}

func SetMotorDir(dirn elev_motor_direction_t) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(dirn))
}

func GetFloorSensor() int {
	return int(C.elev_get_floor_sensor_signal())
}

func GetButtonSignal(button elev_button_type_t, floor int) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor)))
}

func SetFloorIndicatorLamp(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func SetButtonLamp(button elev_button_type_t, floor int, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}

func SetDoorOpenLamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}