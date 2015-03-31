package main

import (
	"Heisprosjekt/io"
	"Heisprosjekt/src"
	"time"
)


func main() {
	command1 := src.Command{src.SET_MOTOR_DIR,src.DIR_UP,src.FLOOR_NONE,src.BUTTON_NONE}
	command2 := src.Command{src.SET_MOTOR_DIR,src.DIR_DOWN,src.FLOOR_NONE,src.BUTTON_NONE}
	command3 := src.Command{src.SET_MOTOR_DIR,src.DIR_STOP,src.FLOOR_NONE,src.BUTTON_NONE}

	chCommandFromControl 	:= make(chan src.Command,1000)
	chButtonOrderToControl 	:= make(chan src.ButtonOrder,1000)
	chFloorSensorToControl	:= make(chan int,1000)

	io.InitIo(chCommandFromControl, chButtonOrderToControl,chFloorSensorToControl)

	chCommandFromControl <- command1
	time.Sleep(1*time.Second)
	chCommandFromControl <- command2
	time.Sleep(1*time.Second)
	chCommandFromControl <- command3
	time.Sleep(1*time.Second)
}