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

	chFloorToQueue				:= make(chan int,2)
	chOrderToQueue 				:= make(chan src.ButtonOrder, 10)
	//chDirectionToQueue			:= make(chan int, 10)
	chOrderFinishedToQueue 		:= make(chan bool,2)
	chAllOrdersFromQueue 		:= make(chan [src.N_FLOORS][3]int,2)
	chDestinationFloorFromQueue	:= make(chan int,2)
	
	chStartTimer 				:= make(chan bool)
	chTimeOut					:= make(chan bool)


	io.InitIo(chCommandToIo, chOrderFromIo, chFloorFromIo)
	goUpUntilReachFloor(chCommandToIo, chFloorFromIo)	
	
	state 			= IDLE
	isOrderFinished = true
	
	go queue.InitQueue(chFloorToQueue, chOrderToQueue, chOrderFinishedToQueue, chAllOrdersFromQueue, chDestinationFloorFromQueue)
	chFloorToQueue <- currentFloor

	go controllerManager(chCommandToIo, chOrderFromIo, chFloorFromIo, chFloorToQueue, 
						 chOrderToQueue, chOrderFinishedToQueue, 
						 chAllOrdersFromQueue, chDestinationFloorFromQueue,chStartTimer,chTimeOut)
	
	go doorTimer(chStartTimer, chTimeOut)
}

func controllerManager( chCommandToIo chan io.Command,
					 	chOrderFromIo chan src.ButtonOrder,
					 	chFloorFromIo chan int, 
					 	chFloorToQueue chan int, 
					 	chOrderToQueue chan src.ButtonOrder,
					 	chOrderFinishedToQueue chan bool, 
					 	chAllOrdersFromQueue chan [src.N_FLOORS][3]int,
					 	chDestinationFloorFromQueue chan int, 
					 	chStartTimer chan bool, 
					 	chTimeOut chan bool){
	for {
		select {
			
			case orderFromIo := <- chOrderFromIo:
				 chOrderToQueue <- orderFromIo

			case currentFloor = <-chFloorFromIo:
				 chCommandToIo  <- io.Command{io.SET_FLOOR_INDICATOR_LAMP,src.ON, currentFloor, src.BUTTON_NONE}
				 chFloorToQueue <- currentFloor
				 destinationFloor = <- chDestinationFloorFromQueue
				
				switch state{
					case MOVING:
						elevatorDirection := chooseElevatorDirection()
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
						//chDirectionToQueue <- elevatorDirection
						
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

					// default:
					// 	time.Sleep(10*time.Millisecond)
				}							
		}
	}
}

func setLights(knowOrders [src.N_FLOORS][3]int, chCommandToIo chan io.Command){
	for floor := 0; floor < src.N_FLOORS; floor ++ {
		chCommandToIo <- io.Command{io.SET_BUTTON_LAMP,knowOrders[floor][src.BUTTON_UP],     floor,src.BUTTON_UP}
		chCommandToIo <- io.Command{io.SET_BUTTON_LAMP,knowOrders[floor][src.BUTTON_DOWN],   floor,src.BUTTON_DOWN}
		chCommandToIo <- io.Command{io.SET_BUTTON_LAMP,knowOrders[floor][src.BUTTON_INSIDE], floor,src.BUTTON_INSIDE}
	}
}

func goUpUntilReachFloor(chCommandToIo chan io.Command, chFloorFromIo chan int){
	chCommandToIo  	<- io.Command{io.SET_MOTOR_DIR,src.DIR_UP,-1,src.BUTTON_NONE}
	currentFloor	= <- chFloorFromIo
	chCommandToIo  	<- io.Command{io.SET_MOTOR_DIR,src.DIR_STOP,-1,src.BUTTON_NONE}
	chCommandToIo  	<- io.Command{io.SET_FLOOR_INDICATOR_LAMP,src.ON,currentFloor,src.BUTTON_NONE}
}

func doorTimer(chStartTimer chan bool, chTimeOut chan bool){
	
	for{
		<-chStartTimer
		time.Sleep(2*time.Second)
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