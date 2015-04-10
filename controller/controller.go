package controller

import (
	"Heisprosjekt/io"
	//"Heisprosjekt/queue"
	"Heisprosjekt/src"
	"Heisprosjekt/tools"
)

const (
	IDLE = 0
	MOVING = 1
	DOOR_OPEN = 2
)

var state int 
var nextFloor int
var currentFloor int

var chCommandFromControl 	= make(chan src.Command)
var chButtonOrderToControl 	= make(chan src.ButtonOrder)
var chFloorSensorToControl	= make(chan int)


var chNewFloor				= make(chan int)
var chNewOrder 				= make(chan src.ButtonOrder)
var chNewDirection			= make(chan int)
var chOrderIsFinished 		= make(chan src.ButtonOrder)
var chNewOrdersFromQueue 	= make(chan src.ElevatorData)
var chNewNextFloorFromQueue = make(chan int)

func InitController() {
	
	io.InitIo(chCommandFromControl, chButtonOrderToControl, chFloorSensorToControl)
	
	goDownUntilReachFloor()	
	state = IDLE
	// InitQueue(channels here...) // Let Queue init network
	
	go controllerHandler()
}

func controllerHandler(){
	for {
		select {
			// case ioOrder:= <- chButtonOrderToControl:
			// 	switch state{
			// 		case IDLE:

			// 		case MOVING:


			// 		case DOOR_OPEN:

			// 	}
			// 	chNewOrder <- ioOrder

			// case currentFloor = <-chFloorSensorToControl:
			// 	switch state{
			// 		case IDLE:
			// 			if currentFloor != nextFloor{
			// 				//move to the next floor
			// 			}
			// 		case MOVING:
			// 			if currentFloor != nextFloor{
			// 				//move to the next floor
			// 			}
			// 		default:
			// 			continue
			// 	}
			// 	chNewFloor <- currentFloor

			case queueOrders := <-chNewOrdersFromQueue:
				// switch state{
				// 	default:
				println("hei")
				setLights(queueOrders)
				// }

			// case nextFloor = <-chNewNextFloorFromQueue:
			// 	continue				
		}
	}
}

func setLights(knowOrders src.ElevatorData){
	println("Enter setLights!")
	tools.PrintQueue(knowOrders) // DELETE THIS!!
	for floor := 0; floor < src.N_FLOORS; floor ++ {
		println(floor)
		println(knowOrders.OutsideOrders[floor][src.BUTTON_UP])

		chCommandFromControl <- src.Command{src.SET_BUTTON_LAMP,knowOrders.OutsideOrders[floor][src.BUTTON_UP],     floor,src.BUTTON_UP}
		chCommandFromControl <- src.Command{src.SET_BUTTON_LAMP,knowOrders.OutsideOrders[floor][src.BUTTON_DOWN],   floor,src.BUTTON_DOWN}
		chCommandFromControl <- src.Command{src.SET_BUTTON_LAMP,knowOrders.InsideOrders[floor], 					floor,src.BUTTON_INSIDE}
	}
	println("Goodbye, end of setLights!")
}

func goDownUntilReachFloor(){

	chCommandFromControl <- src.Command{src.SET_MOTOR_DIR,src.DIR_DOWN,-1,src.BUTTON_NONE}
	currentFloor = <- chFloorSensorToControl
	chCommandFromControl <- src.Command{src.SET_MOTOR_DIR,src.DIR_STOP,-1,src.BUTTON_NONE}
	chCommandFromControl <- src.Command{src.SET_FLOOR_INDICATOR_LAMP,src.ON,currentFloor,src.BUTTON_NONE}

}