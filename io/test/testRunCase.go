package main

import (
	"Heisprosjekt/io"
	"Heisprosjekt/src"
	"fmt"
)

var floor = 3


func initControl(){
	chCommandFromControl 	:= make(chan src.Command)
	chButtonOrderToControl 	:= make(chan src.ButtonOrder)
	chFloorSensorToControl	:= make(chan int)
	io.InitIo(chCommandFromControl, chButtonOrderToControl,chFloorSensorToControl)
	controlHandler(chCommandFromControl,chButtonOrderToControl,chFloorSensorToControl);
}

func main() {
	initControl();
}

func controlHandler(chCommandFromControl chan src.Command, chButtonOrderToControl chan src.ButtonOrder,chFloorSensorToControl chan int){
	fmt.Println("Inside controlHandler.....")
	for{
		select{
			case order := <- chButtonOrderToControl:
				floor = order.Floor
				fmt.Println("Got order sending in command to motor.....")
				commandOrder := src.Command{src.SET_MOTOR_DIR,src.DIR_UP,src.FLOOR_NONE,src.BUTTON_NONE}
				chCommandFromControl <- commandOrder
			case floorArrived := <- chFloorSensorToControl:
				if(floor==floorArrived){
					fmt.Println("Arrived at floor. Stoping...")
					commandFloor := src.Command{src.SET_MOTOR_DIR,src.DIR_STOP,src.FLOOR_NONE,src.BUTTON_NONE}
					chCommandFromControl <- commandFloor				
				}
		}
	}
}
