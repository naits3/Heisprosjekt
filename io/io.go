package io
/*
#cgo CFLAGS: -std=c99
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "channels.h"
#include "elev.h"
*/
import "C"
import "Heisprosjekt/src"
import "fmt"
import "time"


type Command struct{
	CommandType int
	SetValue int
	Floor int
	ButtonType int
}

const(
	SET_BUTTON_LAMP				= 0
	SET_MOTOR_DIR				= 1
	SET_FLOOR_INDICATOR_LAMP	= 2
	SET_DOOR_OPEN_LAMP			= 3
)



func InitIo(chCommandFromControl chan Command, chOrderToControl chan src.ButtonOrder, chFloorToControl chan int){

	err := C.elev_init()
	if  err < 0 { fmt.Println("Could not initialize hardware") }

	var chOrder = make(chan src.ButtonOrder)
	var chFloor = make(chan int)
	
	go ioManager(chCommandFromControl,chOrderToControl,chFloorToControl,chOrder,chFloor)
	go pollFloor(chFloor)
	go pollOrder(chOrder)
}

func ioManager(	chCommandFromControl chan Command, 
				chOrderToControl 	 chan src.ButtonOrder,
			   	chFloorToControl 	 chan int, 
			   	chOrder 		     chan src.ButtonOrder, 
			   	chFloor 			 chan int){
	
	for{
		select{
			case order :=<- chOrder:
				chOrderToControl <- order

			case floor := <- chFloor:
				chFloorToControl <- floor

			case command := <- chCommandFromControl:
				runCommand(command)
		}
	}
}

func pollOrder(chOrder chan src.ButtonOrder){
	for{
		for floor := 0; floor < src.N_FLOORS; floor++{
			
			if(floor < src.N_FLOORS-1){
				if (C.elev_get_button_signal(src.BUTTON_UP, C.int(floor)) == 1) {
					chOrder <- src.ButtonOrder{floor, src.BUTTON_UP}
				}
			}
			
			if(floor > 0) { 
				if (C.elev_get_button_signal(src.BUTTON_DOWN, C.int(floor)) == 1) {
					chOrder <- src.ButtonOrder{floor, src.BUTTON_DOWN}
				}
			}
			
			if (C.elev_get_button_signal(src.BUTTON_INSIDE, C.int(floor)) == 1) {
				chOrder <- src.ButtonOrder{floor, src.BUTTON_INSIDE}
			}
		}
		time.Sleep(40*time.Millisecond)
	}
}

func pollFloor(chFloor chan int){
	previousFloor := -1

	for {
		floor := int(C.elev_get_floor_sensor_signal())
		if floor != -1 && floor != previousFloor{
				chFloor <- floor	
				previousFloor = floor
		}
		time.Sleep(100*time.Millisecond)
	}
}

func runCommand(command Command){
	
	switch command.CommandType{
		
		case SET_MOTOR_DIR:
			C.elev_set_motor_direction(C.elev_motor_direction_t(command.SetValue))
		
		case SET_BUTTON_LAMP:
			switch command.ButtonType{
				
				case src.BUTTON_UP:
					if(command.Floor < src.N_FLOORS-1) {
						C.elev_set_button_lamp(C.elev_button_type_t(command.ButtonType), C.int(command.Floor), C.int(command.SetValue))
					}
				
				case src.BUTTON_DOWN:
					if(command.Floor > 0){
						C.elev_set_button_lamp(C.elev_button_type_t(command.ButtonType), C.int(command.Floor), C.int(command.SetValue))
					}

				case src.BUTTON_INSIDE:
					C.elev_set_button_lamp(C.elev_button_type_t(command.ButtonType), C.int(command.Floor), C.int(command.SetValue))
			
			}		
			

		case SET_FLOOR_INDICATOR_LAMP:
			C.elev_set_floor_indicator(C.int(command.Floor))
		
		case SET_DOOR_OPEN_LAMP:
			C.elev_set_door_open_lamp(C.int(command.SetValue))
	}
}

