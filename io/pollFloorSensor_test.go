package io

import (
	"Heisprosjekt/io"
	"Heisprosjekt/src"
	"testing"
	// "time"
	"fmt"
)

func TestPollFloorSensor(t *testing.T){

	// command1 := src.Command{src.SET_MOTOR_DIR,src.DIR_UP,src.FLOOR_NONE,src.BUTTON_NONE}
	// command2 := src.Command{src.SET_MOTOR_DIR,src.DIR_STOP,src.FLOOR_NONE,src.BUTTON_NONE}
	
	chCommandFromControl 	:= make(chan src.Command)
	chButtonOrderToControl 	:= make(chan src.ButtonOrder)
	chFloorSensorToControl	:= make(chan int)

	io.InitIo(chCommandFromControl, chButtonOrderToControl,chFloorSensorToControl)

	// chCommandFromControl <- command1

	floor := <- chFloorSensorToControl
	fmt.Println(floor)

	u := 1
	if u != 1{
		t.Errorf("Error have happend")
	}
}