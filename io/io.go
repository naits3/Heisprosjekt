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


func InitializeIo(chCommandFromControl chan src.Command, chButtonOrderToControl chan src.ButtonOrder, chFloorSensorToControl chan int){

	chButtonOrder := make(chan src.ButtonOrder)
	chFloorSensor  	:= make(chan int)

	runtime.GOMAXPROCS(runtime.NumCPU())

	if err:=C.elev_init(); err<0{
		fmt.Println("Could not initialize hardware")
		//TODO Handle error, exit. Try one more time
	}

	go ioHandler(chCommandFromControl,chButtonOrderToControl,chFloorSensorToControl,chButtonOrder,chFloorSensor)

	// go pollButtonOrders(chButtonOrder)

	// go pollFloorSensors(chFloorSensor)
}


func ioHandler(chCommandFromControl chan src.Command, chButtonOrderToControl chan src.ButtonOrder,
					chFloorSensorToControl chan int, chButtonOrder chan src.ButtonOrder, chFloorSensor chan int){
	for{
		select{
			case order:=<- chButtonOrder:
				//send order
				chButtonOrderToControl <- order

			case floor:= <- chFloorSensor:
				//send newfloor
				chFloorSensorToControl <- floor

			case command := <- chCommandFromControl:
				doCommand(command)
		}
	}
}

// func pollButtonOrders(chButtonOrder chan src.ButtonOrder){
// 	//TODO: Test funksjonen og fiks problemet med slices!!!
// 	for{
// 		for floor := 0; floor < n_floor; floor++ {
			
// 			if(floor > 0 && driver.GetButtonSignal(src.BUTTON_UP,floor)==1){
// 				chButtonOrder <- []int{0,floor}
// 			}
			
// 			if(floor < n_floor-2 && driver.GetButtonSignal(src.BUTTON_DOWN,floor)==1){
// 				chButtonOrder <- []int{1,floor}
// 			}
			
// 			if(driver.GetButtonSignal(src.BUTTON_INSIDE,floor)==1){
// 				chButtonOrder <- []int{2,floor}
// 			}
// 		}
// 	}
// }

// func pollFloorSensors(chFloorSensor chan int){
// 	//TODO: Test funksjonen
// 	for{	
// 		if floor := driver.GetFloorSensor(); floor != -1{
// 		chFloorSensor <- floor
// 		}
// 	}
// }


func doCommand(command src.Command){
	//TODO  Fiks enum problemet med typer.

	switch commandType := command.CommandType; commandType{
		
		case src.SET_MOTOR_DIR:
			fmt.Println(command.Value)
			C.elev_set_motor_direction(C.elev_motor_direction_t(command.Value))
		
		case src.SET_BUTTON_LAMP:
			//SET BUTTON LAMP
			// driver.SetButtonLamp(1, command.floor, command.value)
		
		case src.SET_FLOOR_INDICATOR_LAMP:
			// SET FLOOR INDICATOR
			// driver.SetFloorIndicatorLamp(command.value)
		
		case src.SET_DOOR_OPEN_LAMP:
			//Set door open lamp
	}

}

