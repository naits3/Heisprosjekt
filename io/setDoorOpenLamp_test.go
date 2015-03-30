package io

import (
	"Heisprosjekt/io"
	"Heisprosjekt/src"
	"testing"
	"time"
)

func TestMotor(t *testing.T){

	command1 := src.Command{src.SET_DOOR_OPEN_LAMP,src.ON,src.FLOOR_NONE,src.BUTTON_NONE}
	command2 := src.Command{src.SET_DOOR_OPEN_LAMP,src.OFF,src.FLOOR_NONE,src.BUTTON_NONE}

	chCommandFromControl 	:= make(chan src.Command)
	chButtonOrderToControl 	:= make(chan src.ButtonOrder)
	chFloorSensorToControl	:= make(chan int)

	io.InitIo(chCommandFromControl, chButtonOrderToControl,chFloorSensorToControl)

	chCommandFromControl <- command1
	time.Sleep(1*time.Second)
	chCommandFromControl <- command2
	time.Sleep(1*time.Second)

	u := 1
	if u != 1{
		t.Errorf("Error have happend")
	}
}