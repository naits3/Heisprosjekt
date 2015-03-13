package io

import "Heisprosjekt/driver"
import "Heisprosjekt/source"
import "runtime"
import "fmt"

type Commmad struct{
	commandType int
	floor 		int
	value 		int
	buttonType 	int
}

const(
	SET_BUTTON_LAMP				= 0
	SET_MOTOR_DIR				= 1
	SET_FLOOR_INDICATOR_LAMP	= 2
)

const(
	OFF			=  0
	ON 			=  1
)


var n_floor

func InitializeIo(chCommand chan Command, chButtonOrderToControl chan source.ButtonOrder){

	runtime.GOMAXPROCS(runtime.NumCPU())

	n_floor = source.GetNFloors()

	chButtonOrder := make(chan source.BtnOrder)
	chFloorSensor  	:= make(chan int)
	
	if err:=driver.Init(); err<0{
		fmt.Println("Could not initialize hardware")
		//TODO Handle error, exit. Try one more time
	}

	go ioHandler(chCommand,chButtonOrderToControl,chButtonOrder,chFloorSensor)

	go pollButtonOrders(chButtonOrder)

	go pollFloorSensors(chFloorSensor)
}

func ioHandler(chCommandFromControl chan Command, chButtonOrderToControl chan source.ButtonOrder,
					chFloorSensorToControl chan int, chButtonOrder chan source.ButtonOrder, chFloorSensor chan int){
	for{
		select{
			case order:=<- chButtonOrder:
				//send order
				chButtonOrderToControl <- 1

			case floor:= <- chFloorSensor:
				//send newfloor
				chButtonOrderToControl <- 2

			case command := <- chCommand:
				doCommand(command)
		}
	}
}

func pollButtonOrders(chButtonOrder chan source.ButtonOrder){
	//TODO: Test funksjonen og fiks problemet med slices!!!
	for{
		for floor := 0; floor < n_floor; floor++ {
			
			if(floor > 0 && driver.GetButtonSignal(driver.BUTTON_CALL_UP,floor)==1){
				chButtonOrder <- []int{0,floor}
			}
			
			if(floor < n_floor-2 && driver.GetButtonSignal(driver.BUTTON_CALL_DOWN,floor)==1){
				chButtonOrder <- []int{1,floor}
			}
			
			if(driver.GetButtonSignal(driver.BUTTON_COMMAND,floor)==1){
				chButtonOrder <- []int{2,floor}
			}
		}
	}
}

func pollFloorSensors(chFloorSensor chan int){
	//TODO: Test funksjonen
	for{	
		if floor := driver.GetFloorSensor(); floor != -1{
		chFloorSensor <- floor
		}
	}
}

func doCommand(command Command){
	//TODO  Fiks enum problemet med typer.

	if(SET_BUTTON_LAMP == command.CommandType){
		driver.SetButtonLamp(1, CommandControl.floor, CommandControl.value)
	}

	if(SET_MOTOR_DIR == command.CommandType){
		driver.SetMotorDir(1)
	}

	if(SET_FLOOR_INDICATOR_LAMP == command.CommandType){
		driver.SetFloorIndicatorLamp(2)
	}
	//implement SetDoorOpenLamp
}

