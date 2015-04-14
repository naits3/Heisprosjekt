package controller

import (
	"Heisprosjekt/io"
	"Heisprosjekt/queue"
	"Heisprosjekt/src"
	"Heisprosjekt/tools"
	"time"
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

	queue.InitQueue(channels here...) // Let Queue init network
	
	go controllerHandler()
}

func controllerHandler(){
	for {
		select {
			case ioOrder:= <- chButtonOrderToControl:
				chNewOrder <- ioOrder

			case currentFloor = <-chFloorSensorToControl:
				switch state{
					case IDLE:
						if currentFloor != nextFloor{
							dir := findElevatorDirection()
							chCommandFromControl <- src.Command{src.SET_MOTOR_DIR,dir,-1,src.BUTTON_NONE}
							state = MOVING
						}else{
							chCommandFromControl <- src.Command{src.SET_MOTOR_DIR,src.DIR_STOP,-1,src.BUTTON_NONE}
							state = IDLE
							//Hvordan skal vi vite at det ikke er en inside ordere??
						}
					case MOVING:
						if currentFloor != nextFloor{
							dir := findElevatorDirection()
							chCommandFromControl <- src.Command{src.SET_MOTOR_DIR,dir,-1,src.BUTTON_NONE}
							state = MOVING
						}else{
							chCommandFromControl <- src.Command{src.SET_MOTOR_DIR,src.DIR_STOP,-1,src.BUTTON_NONE}
							state = DOOR_OPEN
							go openDoor()
						}
					case DOOR_OPEN:
						continue
				}
				//chNewFloor <- currentFloor

			case queueOrders := <-chNewOrdersFromQueue:
				setLights(queueOrders)
		

			case nextFloor = <-chNewNextFloorFromQueue:
				continue
		}				
	}
}

func setLights(knowOrders src.ElevatorData){
	//println("Enter setLights!")
	tools.PrintQueue(knowOrders) // DELETE THIS!!
	for floor := 0; floor < src.N_FLOORS; floor ++ {
		//println(floor)
		println(knowOrders.OutsideOrders[floor][src.BUTTON_UP])

		chCommandFromControl <- src.Command{src.SET_BUTTON_LAMP,knowOrders.OutsideOrders[floor][src.BUTTON_UP],     floor,src.BUTTON_UP}
		chCommandFromControl <- src.Command{src.SET_BUTTON_LAMP,knowOrders.OutsideOrders[floor][src.BUTTON_DOWN],   floor,src.BUTTON_DOWN}
		chCommandFromControl <- src.Command{src.SET_BUTTON_LAMP,knowOrders.InsideOrders[floor], 					floor,src.BUTTON_INSIDE}
	}
	//println("Goodbye, end of setLights!")
}

func goDownUntilReachFloor(){
	chCommandFromControl <- src.Command{src.SET_MOTOR_DIR,src.DIR_DOWN,-1,src.BUTTON_NONE}
	
	currentFloor	= <- chFloorSensorToControl
	nextFloor 		= currentFloor
	
	chCommandFromControl <- src.Command{src.SET_MOTOR_DIR,src.DIR_STOP,-1,src.BUTTON_NONE}
	chCommandFromControl <- src.Command{src.SET_FLOOR_INDICATOR_LAMP,src.ON,currentFloor,src.BUTTON_NONE}
}


func openDoor(){
    command1 := src.Command{src.SET_DOOR_OPEN_LAMP,src.ON,-1,src.BUTTON_NONE}
	command2 := src.Command{src.SET_DOOR_OPEN_LAMP,src.OFF,-1,src.BUTTON_NONE}
	chCommandFromControl <- command1
	time.Sleep(3*time.Second)
	chCommandFromControl <- command2
	state = IDLE
}

func findElevatorDirection() int {
	if(currentFloor<nextFloor){
		return src.DIR_UP
	}else{
		return src.DIR_DOWN
	}
}