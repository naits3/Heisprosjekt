package io

import (
	"Heisprosjekt/io"
	"Heisprosjekt/src"
	"testing"
	"time"
)


func TestFloorIndicator(t *testing.T){

	command1 := src.Command{src.SET_FLOOR_INDICATOR_LAMP,0,src.FLOOR_2,src.BUTTON_NONE}
	command2 := src.Command{src.SET_FLOOR_INDICATOR_LAMP,0,src.FLOOR_3,src.BUTTON_NONE}
	command3 := src.Command{src.SET_FLOOR_INDICATOR_LAMP,0,src.FLOOR_4,src.BUTTON_NONE}

	chCommandFromControl 	:= make(chan src.Command)
	chButtonOrderToControl 	:= make(chan src.ButtonOrder)
	chFloorSensorToControl	:= make(chan int)

	io.InitIo(chCommandFromControl, chButtonOrderToControl,chFloorSensorToControl)

	chCommandFromControl <- command1
	time.Sleep(1*time.Second)
	chCommandFromControl <- command2
	time.Sleep(1*time.Second)
	chCommandFromControl <- command3
	time.Sleep(1*time.Second)


	u := 1
	if u != 1{
		t.Errorf("Error have happend")
	}
}


