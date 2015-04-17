package controller

import (
	"Heisprosjekt/io"
	"Heisprosjekt/queue"
	"Heisprosjekt/src"
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


func InitController() {
	
	chCommandToIo 				:= make(chan io.Command)
	chOrderFromIo 				:= make(chan src.ButtonOrder)
	chFloorFromIo				:= make(chan int)

	chFloorToQueue				:= make(chan int)
	chOrderToQueue 				:= make(chan src.ButtonOrder)
	chDirectionToQueue			:= make(chan int)
	chOrderFinishedToQueue 		:= make(chan bool)
	chAllOrdersFromQueue 		:= make(chan src.ElevatorData)
	chDestinationFloorFromQueue	:= make(chan int)
	
	chStartTimer 				:= make(chan bool)
	chTimeOut					:= make(chan bool)


	io.InitIo(chCommandToIo, chOrderFromIo, chFloorFromIo)
	goDownUntilReachFloor(chCommandToIo, chFloorFromIo)	
	
	state 			= IDLE
	isOrderFinished = true
	
	queue.InitQueue(chFloorToQueue, chOrderToQueue, chDirectionToQueue, chOrderFinishedToQueue, 
					chAllOrdersFromQueue, chDestinationFloorFromQueue)

	chFloorToQueue <- currentFloor

	go controllerManager(chCommandToIo, chOrderFromIo, chFloorFromIo, chFloorToQueue, 
						 chOrderToQueue, chDirectionToQueue, chOrderFinishedToQueue, 
						 chAllOrdersFromQueue, chDestinationFloorFromQueue,chStartTimer,chTimeOut)
	
	go doorTimer(chStartTimer, chTimeOut)
}

func controllerManager( chCommandToIo chan io.Command,
					 	chOrderFromIo chan src.ButtonOrder,
					 	chFloorFromIo chan int, 
					 	chFloorToQueue chan int, 
					 	chOrderToQueue chan src.ButtonOrder,
					 	chDirectionToQueue chan int,
					 	chOrderFinishedToQueue chan bool, 
					 	chAllOrdersFromQueue chan src.ElevatorData,
					 	chDestinationFloorFromQueue chan int, 
					 	chStartTimer chan bool, 
					 	chTimeOut chan bool){
	for {
		select {
			
			case orderFromIo := <- chOrderFromIo:
				 chOrderToQueue <- orderFromIo

			case currentFloor = <-chFloorFromIo:
				 chCommandToIo  <- io.Command{io.SET_FLOOR_INDICATOR_LAMP,src.ON, currentFloor, src.BUTTON_NONE}
				
				switch state{
					
					case MOVING:
						elevatorDirection := chooseElevatorDirection()
						chDirectionToQueue 	<- elevatorDirection
						chCommandToIo 		<- io.Command{io.SET_MOTOR_DIR,elevatorDirection,-1,src.BUTTON_NONE}

						if(currentFloor == destinationFloor){
							state = DOOR_OPEN
							chCommandToIo <- io.Command{io.SET_DOOR_OPEN_LAMP,src.ON,-1,src.BUTTON_NONE}
							chStartTimer  <- true
						}

					default:
						// feilhaandtering?
						continue
				}
				chFloorToQueue <- currentFloor
				
			case <- chTimeOut:
					state = IDLE
					chCommandToIo <- io.Command{io.SET_DOOR_OPEN_LAMP,src.OFF,-1,src.BUTTON_NONE}
					chOrderFinishedToQueue <- true
					isOrderFinished = true

			case ordersFromQueue := <-chAllOrdersFromQueue:
				setLights(ordersFromQueue,chCommandToIo)
		
			case destinationFloor = <- chDestinationFloorFromQueue:
				
				switch state{
					
					case IDLE:
						
						elevatorDirection := chooseElevatorDirection()
						chDirectionToQueue <- elevatorDirection
						
						if(currentFloor == destinationFloor){
							state = DOOR_OPEN
							chCommandToIo <- io.Command{io.SET_DOOR_OPEN_LAMP,src.ON,-1,src.BUTTON_NONE}
							chStartTimer <- true
						
						}else{
							isOrderFinished = false
							chCommandToIo <- io.Command{io.SET_MOTOR_DIR,elevatorDirection,-1,src.BUTTON_NONE}
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

func setLights(knowOrders src.ElevatorData, chCommandToIo chan io.Command){
	for floor := 0; floor < src.N_FLOORS; floor ++ {
		chCommandToIo <- io.Command{io.SET_BUTTON_LAMP,knowOrders.OutsideOrders[floor][src.BUTTON_UP],     floor,src.BUTTON_UP}
		chCommandToIo <- io.Command{io.SET_BUTTON_LAMP,knowOrders.OutsideOrders[floor][src.BUTTON_DOWN],   floor,src.BUTTON_DOWN}
		chCommandToIo <- io.Command{io.SET_BUTTON_LAMP,knowOrders.InsideOrders[floor], 					   floor,src.BUTTON_INSIDE}
	}
}

func goDownUntilReachFloor(chCommandToIo chan io.Command, chFloorFromIo chan int){
	chCommandToIo  	<- io.Command{io.SET_MOTOR_DIR,src.DIR_DOWN,-1,src.BUTTON_NONE}
	currentFloor	= <- chFloorFromIo
	chCommandToIo  	<- io.Command{io.SET_MOTOR_DIR,src.DIR_STOP,-1,src.BUTTON_NONE}
	chCommandToIo  	<- io.Command{io.SET_FLOOR_INDICATOR_LAMP,src.ON,currentFloor,src.BUTTON_NONE}
}

func doorTimer(chStartTimer chan bool, chTimeOut chan bool){
	
	for{
		<-chStartTimer
		time.Sleep(3*time.Second)
		chTimeOut <- true		
	}
}

func chooseElevatorDirection() int {
	
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
Fix floorindicator 																									| OK
Error handling when reach new floor
change name on chAllOrdersFromQueue and queueORder variable 														|		
*/