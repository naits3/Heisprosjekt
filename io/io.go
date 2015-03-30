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
import "runtime"
import "fmt"

func InitIo(chCommandFromControl chan src.Command, chButtonOrderToControl chan src.ButtonOrder, chFloorSensorToControl chan int){

	chButtonOrder := make(chan src.ButtonOrder)
	chFloorSensor  	:= make(chan int)

	runtime.GOMAXPROCS(runtime.NumCPU())

	if err:=C.elev_init(); err<0{
		fmt.Println("Could not initialize hardware")
		//TODO Handle error, exit. Try one more time
	}

	go ioHandler(chCommandFromControl,chButtonOrderToControl,chFloorSensorToControl,chButtonOrder,chFloorSensor)
	go pollFloorSensors(chFloorSensor)
	go pollButtonOrders(chButtonOrder)
}

func ioHandler(chCommandFromControl chan src.Command, chButtonOrderToControl chan src.ButtonOrder,chFloorSensorToControl chan int, chButtonOrder chan src.ButtonOrder, chFloorSensor chan int){
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

func pollButtonOrders(chButtonOrder chan src.ButtonOrder){
	for{
		for floor := 0; floor < src.N_FLOORS; floor++{
			
			if(floor > src.FLOOR_1 && driver.GetButtonSignal(src.BUTTON_UP,floor)==1){
				chButtonOrder <- src.ButtonOrder{floor, src.BUTTON_UP}
			}
			
			if(floor < src.N_FLOOR-2 && driver.GetButtonSignal(src.BUTTON_DOWN,floor)==1){
				chButtonOrder <- src.ButtonOrder{floor, src.BUTTON_DOWN}
			}
			
			if(driver.GetButtonSignal(src.BUTTON_INSIDE,floor)==1){
				chButtonOrder <- src.ButtonOrder{floor, src.BUTTON_INSIDE}
			}
		}
	}
}

func pollFloorSensors(chFloorSensor chan int){
	for{
		if floor := int(C.elev_get_floor_sensor_signal()); floor != -1{
			chFloorSensor <- floor
		}
	}
}

func doCommand(command src.Command){
	//TODO  Fiks enum problemet med typer.

	switch commandType := command.CommandType; commandType{
		
		case src.SET_MOTOR_DIR:
			C.elev_set_motor_direction(C.elev_motor_direction_t(command.SetValue))
		
		case src.SET_BUTTON_LAMP:
			C.elev_set_button_lamp(C.elev_button_type_t(command.ButtonType), C.int(command.Floor), C.int(command.SetValue))
		
		case src.SET_FLOOR_INDICATOR_LAMP:
			C.elev_set_floor_indicator(C.int(command.Floor))
		
		case src.SET_DOOR_OPEN_LAMP:
			C.elev_set_door_open_lamp(C.int(command.SetValue))
	}

}

