package main

import (
	"Heisprosjekt/io"
	"Heisprosjekt/src"
	"fmt"
)


func main() {
	chCommandFromControl 	:= make(chan src.Command)
	chButtonOrderToControl 	:= make(chan src.ButtonOrder)
	chFloorSensorToControl	:= make(chan int)
	io.InitIo(chCommandFromControl, chButtonOrderToControl,chFloorSensorToControl)

	type btnOrder src.ButtonOrder
	for{
		btnOrder := <- chButtonOrderToControl
		fmt.Println(btnOrder.Floor)
		fmt.Println(btnOrder.ButtonType)
	}

	// for{
	// 	a := driver.GetFloorSensor()
	// 	fmt.Println(a)
	// 	time.Sleep(1*time.Second)
	// }
}