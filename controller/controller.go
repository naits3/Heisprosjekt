package controller

import (
	//"Heisprosjekt/io"
	//"Heisprosjekt/queue"
)

const (
	IDLE = 0
	MOVING = 1
	DOOR_OPEN = 2
)

var state int 
var nextFloor int
var currentFloor int

func InitController() {
	// chCommandFromControl 	:= make(chan src.Command)
	// chButtonOrderToControl 	:= make(chan src.ButtonOrder)
	// chFloorSensorToControl	:= make(chan int)


	// chNewFloor		:= make(chan int)
	// chNewOrder 		:= make(chan src.ButtonOrder)
	// chNewDirection	:= make(chan int)
	// chOrderIsFinished 	:= make(chan src.ButtonOrder)
	// chNewOrdersFromQueue := make(chan src.NewStruct)
	// Create queueChannels!

	//io.InitIo(chCommandFromControl, chButtonOrderToControl, chFloorSensorToControl) // anta io sletter lys
	// Go down until you reach a floor
	// Set correct floor light

	state = IDLE
	// InitQueue(channels here...) // Let Queue init network
}

func controllerHandler() {
	for {
		select {
			case ioOrder:= <- chButtonOrderToControl:
				switch state{
					case IDLE:

					case MOVING:

					case DOOR_OPEN:

				}
				chNewOrder <- ioOrder

			case currentFloor = <-chFloorSensorToControl:
				switch state{
					case IDLE:
						if currentFloor != nextFloor{
							//move to the next floor
						}
					case MOVING:
						if currentFloor != nextFloor{
							//move to the next floor
						}
					default:
						continue
				}
				chNewFloor <- currentFloor


			case QueueOrders := <-chNewOrdersFromQueue:
				switch state{
					default:
						setLights(QueueOrders[0])
						nextFloor = QueueOrders[1]
				}				
			}
	}
}

func setLights(knownOrders src.ElevatorData){
	for floor := 1; floor < src.N_FLOORS; floor ++ {
		chCommandFromControl <- src.Command{src.SET_BUTTON_LAMP,knowOrders.OutsideOrders[floor][src.BUTTON_UP],     floor,src.BUTTON_UP}
		chCommandFromControl <- src.Command{src.SET_BUTTON_LAMP,knowOrders.OutsideOrders[floor][src.BUTTON_DOWN],   floor,src.BUTTON_DOWN}
		chCommandFromControl <- src.Command{src.SET_BUTTON_LAMP,knowOrders.InsideOrders[floor], 					floor,src.BUTTON_INSIDE}
	}
}