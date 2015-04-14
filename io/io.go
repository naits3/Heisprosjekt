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


var chButtonOrder = make(chan src.ButtonOrder)
var chFloorSensor = make(chan int)

func InitIo(chCommandFromControl chan src.Command, chButtonOrderToControl chan src.ButtonOrder, chFloorSensorToControl chan int){

	if err:=C.elev_init(); err<0{
		fmt.Println("Could not initialize hardware")
	}

	go ioHandler(chCommandFromControl,chButtonOrderToControl,chFloorSensorToControl)
	go pollFloorSensors()
	go pollButtonOrders()
}

func ioHandler(chCommandFromControl chan src.Command, chButtonOrderToControl chan src.ButtonOrder,chFloorSensorToControl chan int){
	for{
		select{
			case order :=<- chButtonOrder:
				chButtonOrderToControl <- order

			case floor := <- chFloorSensor:
				chFloorSensorToControl <- floor

			case command := <- chCommandFromControl:
				doCommand(command)
		}
	}
}

func pollButtonOrders(){
	for{
		for floor := 0; floor < src.N_FLOORS; floor++{
			
			if(floor < src.N_FLOORS-1){
				if(int(C.elev_get_button_signal(C.elev_button_type_t(src.BUTTON_UP), C.int(floor)))==1){
					chButtonOrder <- src.ButtonOrder{floor, src.BUTTON_UP}
				}
			}
			
			if(floor > 0) { 
				if(int(C.elev_get_button_signal(C.elev_button_type_t(src.BUTTON_DOWN), C.int(floor)))==1){
					chButtonOrder <- src.ButtonOrder{floor, src.BUTTON_DOWN}
				}
			}
			
			if(int(C.elev_get_button_signal(C.elev_button_type_t(src.BUTTON_INSIDE), C.int(floor)))==1){
				chButtonOrder <- src.ButtonOrder{floor, src.BUTTON_INSIDE}
			}
		}
		time.Sleep(40*time.Millisecond)
	}
}

func pollFloorSensors(){
	for{
		if floor := int(C.elev_get_floor_sensor_signal()); floor != -1{
			chFloorSensor <- floor	
		}
		time.Sleep(100*time.Millisecond)
	}
}

func doCommand(command src.Command){
	switch commandType := command.CommandType; commandType{
		case src.SET_MOTOR_DIR:
			C.elev_set_motor_direction(C.elev_motor_direction_t(command.SetValue))
		
		case src.SET_BUTTON_LAMP:
			switch buttonType := command.ButtonType;buttonType{
				case src.BUTTON_UP:
					if(command.Floor < src.N_FLOORS-1) {
						C.elev_set_button_lamp(C.elev_button_type_t(command.ButtonType), C.int(command.Floor), C.int(command.SetValue))
						println(command.SetValue)
					}
				case src.BUTTON_DOWN:
					if(command.Floor > 0){
						C.elev_set_button_lamp(C.elev_button_type_t(command.ButtonType), C.int(command.Floor), C.int(command.SetValue))
						println(command.SetValue)
					}
				case src.BUTTON_INSIDE:
					C.elev_set_button_lamp(C.elev_button_type_t(command.ButtonType), C.int(command.Floor), C.int(command.SetValue))
					println(command.SetValue)
			}		
			

		case src.SET_FLOOR_INDICATOR_LAMP:
			C.elev_set_floor_indicator(C.int(command.Floor))
		
		case src.SET_DOOR_OPEN_LAMP:
			C.elev_set_door_open_lamp(C.int(command.SetValue))
	}
}

