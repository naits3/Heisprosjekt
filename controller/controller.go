package controller

import (
	"Heisprosjekt/io"
	"Heisprosjekt/queue"
	"Heisprosjekt/src"
	//"Heisprosjekt/tools"
	"time"
)

const (
	IDLE = 0
	MOVING = 1
	DOOR_OPEN = 2
)

var state int 
var destinationFloor int
var currentFloor	 int
var isOrderFinished	 bool

var chCommandToIo 			= make(chan src.Command)
var chOrderFromIo 			= make(chan src.ButtonOrder)
var chFloorFromIo			= make(chan int)

var chFloorToQueue			= make(chan int)
var chOrderToQueue 			= make(chan src.ButtonOrder)
var chDirectionToQueue		= make(chan int)
var chOrderFinishedToQueue 	= make(chan bool)
var chAllOrdersFromQueue 	= make(chan src.ElevatorData)
var chDestinationFloorFromQueue	= make(chan int)
var setDoorTimer 			= make(chan bool)
var chTimeOut				= make(chan bool)

func InitController() {
	
	isOrderFinished = true

	io.InitIo(chCommandToIo, chOrderFromIo, chFloorFromIo)
	goDownUntilReachFloor()	
	queue.InitQueue(chFloorToQueue, chOrderToQueue, chDirectionToQueue, chOrderFinishedToQueue, chAllOrdersFromQueue, chDestinationFloorFromQueue) // Let Queue init network
	chFloorToQueue <- currentFloor
	go controllerHandler()
	go doorTimer()
}

func controllerHandler(){
	for {
		select {
			case ioOrder:= <- chOrderFromIo:
				chOrderToQueue <- ioOrder

			case currentFloor = <-chFloorFromIo:
				chCommandToIo  <- src.Command{src.SET_FLOOR_INDICATOR_LAMP,src.ON,currentFloor,src.BUTTON_NONE}
				switch state{
					
					case MOVING:
						motorDirection := findElevatorDirection()
						chDirectionToQueue <- motorDirection
						chCommandToIo <- src.Command{src.SET_MOTOR_DIR,motorDirection,-1,src.BUTTON_NONE}

						if(currentFloor == destinationFloor){
							state = DOOR_OPEN
							chCommandToIo <- src.Command{src.SET_DOOR_OPEN_LAMP,src.ON,-1,src.BUTTON_NONE}
							setDoorTimer <- true
						}

					default:
						// feilhaandtering?
						continue
				}
				chFloorToQueue <- currentFloor
				
			case <- chTimeOut:
					state = IDLE
					chCommandToIo <- src.Command{src.SET_DOOR_OPEN_LAMP,src.OFF,-1,src.BUTTON_NONE}
					chOrderFinishedToQueue <- true
					isOrderFinished = true

			case queueOrders := <-chAllOrdersFromQueue:
				//print("C: received orders from Q")
				setLights(queueOrders)
		
			case destinationFloor = <-  chDestinationFloorFromQueue:
				
				switch state {
					
					case IDLE:
						
						motorDirection := findElevatorDirection()
						chDirectionToQueue <- motorDirection
						
						if(currentFloor == destinationFloor){
							state = DOOR_OPEN
							chCommandToIo <- src.Command{src.SET_DOOR_OPEN_LAMP,src.ON,-1,src.BUTTON_NONE}
							setDoorTimer <- true
						
						}else{
							isOrderFinished = false
							chCommandToIo <- src.Command{src.SET_MOTOR_DIR,motorDirection,-1,src.BUTTON_NONE}
							state = MOVING
						}		

					case MOVING:
						isOrderFinished = false

					default:
						continue
				}							
		}
	}
}

func setLights(knowOrders src.ElevatorData){
	for floor := 0; floor < src.N_FLOORS; floor ++ {
		chCommandToIo <- src.Command{src.SET_BUTTON_LAMP,knowOrders.OutsideOrders[floor][src.BUTTON_UP],     floor,src.BUTTON_UP}
		chCommandToIo <- src.Command{src.SET_BUTTON_LAMP,knowOrders.OutsideOrders[floor][src.BUTTON_DOWN],   floor,src.BUTTON_DOWN}
		chCommandToIo <- src.Command{src.SET_BUTTON_LAMP,knowOrders.InsideOrders[floor], 					floor,src.BUTTON_INSIDE}
	}
}

func goDownUntilReachFloor(){
	
	chCommandToIo  <- src.Command{src.SET_MOTOR_DIR,src.DIR_DOWN,-1,src.BUTTON_NONE}

	currentFloor		 = <- chFloorFromIo
	destinationFloor 	 = currentFloor
	chCommandToIo  <- src.Command{src.SET_MOTOR_DIR,src.DIR_STOP,-1,src.BUTTON_NONE}
	chCommandToIo  <- src.Command{src.SET_FLOOR_INDICATOR_LAMP,src.ON,currentFloor,src.BUTTON_NONE}
	state = IDLE


}


func doorTimer(){

	for{
		<-setDoorTimer
		println("openDoor")
		time.Sleep(3*time.Second)
		println("closeDoor")	
		chTimeOut <- true		
	}
}


func findElevatorDirection() int {
	if(currentFloor<destinationFloor){
		return src.DIR_UP
	}else if(currentFloor>destinationFloor){
		return src.DIR_DOWN
	}else{
		return src.DIR_STOP
	}
}

/*
BUGS
If a ordere for next floor is recived and the next floor equals the current floor, the confirmation wont be sent. 	| OK
We may get some problems with openDoor. 																			|
Fix floorindicator lamp

*/