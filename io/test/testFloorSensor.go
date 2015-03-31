package main

import (
	"Heisprosjekt/io"
	//"Heisprosjekt/driver"
	//"time"
	"Heisprosjekt/src"
	"fmt"
)


func main() {
	chCommandFromControl 	:= make(chan src.Command)
	chButtonOrderToControl 	:= make(chan src.ButtonOrder)
	chFloorSensorToControl	:= make(chan int)
	io.InitIo(chCommandFromControl, chButtonOrderToControl,chFloorSensorToControl)


	for{
		floor := <- chFloorSensorToControl
		fmt.Println(floor)
	}

	// for{
	// 	a := driver.GetFloorSensor()
	// 	fmt.Println(a)
	// 	time.Sleep(1*time.Second)
	// }
}