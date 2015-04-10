package io

import (
	"Heisprosjekt/src"
	"testing"
	"time"
)


func TestMotor(t *testing.T){

	command1 := src.Command{src.SET_MOTOR_DIR,src.DIR_UP,-1,src.BUTTON_NONE}
	command2 := src.Command{src.SET_MOTOR_DIR,src.DIR_DOWN,-1,src.BUTTON_NONE}
	command3 := src.Command{src.SET_MOTOR_DIR,src.DIR_STOP,-1,src.BUTTON_NONE}

	chCommandFromControl 	:= make(chan src.Command)
	chButtonOrderToControl 	:= make(chan src.ButtonOrder)
	chFloorSensorToControl	:= make(chan int)

	InitIo(chCommandFromControl, chButtonOrderToControl,chFloorSensorToControl)

	chCommandFromControl <- command1
	time.Sleep(1*time.Second)
	chCommandFromControl <- command2
	time.Sleep(1*time.Second)
	chCommandFromControl <- command3
	time.Sleep(1*time.Second)
}


func TestSetButtonLight(t *testing.T){

	command1 := src.Command{src.SET_BUTTON_LAMP,src.ON,1,src.BUTTON_INSIDE}
	command2 := src.Command{src.SET_BUTTON_LAMP,src.ON,2,src.BUTTON_UP}
	command3 := src.Command{src.SET_BUTTON_LAMP,src.ON,3,src.BUTTON_DOWN}

	chCommandFromControl 	:= make(chan src.Command)
	chButtonOrderToControl 	:= make(chan src.ButtonOrder)
	chFloorSensorToControl	:= make(chan int)

	InitIo(chCommandFromControl, chButtonOrderToControl,chFloorSensorToControl)

	time.Sleep(1*time.Second)
	chCommandFromControl <- command1
	time.Sleep(1*time.Second)
	chCommandFromControl <- command2
	time.Sleep(1*time.Second)
	chCommandFromControl <- command3
	time.Sleep(1*time.Second)
}


func TestSetDoorOpenLamp(t *testing.T){

	command1 := src.Command{src.SET_DOOR_OPEN_LAMP,src.ON,-1,src.BUTTON_NONE}
	command2 := src.Command{src.SET_DOOR_OPEN_LAMP,src.OFF,-1,src.BUTTON_NONE}

	chCommandFromControl 	:= make(chan src.Command)
	chButtonOrderToControl 	:= make(chan src.ButtonOrder)
	chFloorSensorToControl	:= make(chan int)

	InitIo(chCommandFromControl, chButtonOrderToControl,chFloorSensorToControl)

	chCommandFromControl <- command1
	time.Sleep(1*time.Second)
	chCommandFromControl <- command2
	time.Sleep(1*time.Second)
}


func TestSetFloorIndicator(t *testing.T){

	command1 := src.Command{src.SET_FLOOR_INDICATOR_LAMP,0,1,src.BUTTON_NONE}
	command2 := src.Command{src.SET_FLOOR_INDICATOR_LAMP,0,2,src.BUTTON_NONE}
	command3 := src.Command{src.SET_FLOOR_INDICATOR_LAMP,0,3,src.BUTTON_NONE}

	chCommandFromControl 	:= make(chan src.Command)
	chButtonOrderToControl 	:= make(chan src.ButtonOrder)
	chFloorSensorToControl	:= make(chan int)

	InitIo(chCommandFromControl, chButtonOrderToControl,chFloorSensorToControl)

	chCommandFromControl <- command1
	time.Sleep(1*time.Second)
	chCommandFromControl <- command2
	time.Sleep(1*time.Second)
	chCommandFromControl <- command3
	time.Sleep(1*time.Second)
}


func TestRunCase(t *testing.T){
	
	floor := 0

	chCommandFromControl 	:= make(chan src.Command)
	chButtonOrderToControl 	:= make(chan src.ButtonOrder)
	chFloorSensorToControl	:= make(chan int)
	InitIo(chCommandFromControl, chButtonOrderToControl,chFloorSensorToControl)
	println("Inside controlHandler.....")
		for{
			select{
				case order := <- chButtonOrderToControl:
					floor = order.Floor
					println("Got order sending in command to motor.....")
					commandOrder := src.Command{src.SET_MOTOR_DIR,src.DIR_UP,-1,src.BUTTON_NONE}
					chCommandFromControl <- commandOrder
				case floorArrived := <- chFloorSensorToControl:
					if(floor==floorArrived){
						println("Arrived at floor. Stoping...")
						commandFloor := src.Command{src.SET_MOTOR_DIR,src.DIR_STOP,-1,src.BUTTON_NONE}
						chCommandFromControl <- commandFloor
						break				
					}
			}
		}
}

func TestFloorSensor(t *testing.T){
	chCommandFromControl 	:= make(chan src.Command)
	chButtonOrderToControl 	:= make(chan src.ButtonOrder)
	chFloorSensorToControl	:= make(chan int)
	InitIo(chCommandFromControl, chButtonOrderToControl,chFloorSensorToControl)


	for{
		floor := <- chFloorSensorToControl
		println(floor)
	}
}


func TestPoolButtonOrders(t *testing.T){

	chCommandFromControl 	:= make(chan src.Command)
	chButtonOrderToControl 	:= make(chan src.ButtonOrder)
	chFloorSensorToControl	:= make(chan int)
	InitIo(chCommandFromControl, chButtonOrderToControl,chFloorSensorToControl)

	type btnOrder src.ButtonOrder
	for{
		btnOrder := <- chButtonOrderToControl
		println(btnOrder.Floor)
		println(btnOrder.ButtonType)
	}

}