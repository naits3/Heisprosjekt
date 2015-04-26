package controller

import (
	"Heisprosjekt/io"
	"Heisprosjekt/queue"
	"Heisprosjekt/src"
	"time"
)

/* Name-conventions:
FI = From IO
TI = To IO
FQ = From Queue
TQ = To Queue
*/

const (
	IDLE = 0
	MOVING = 1
	DOOR_OPEN = 2
)

var state 				int 
var destinationFloor 	int
var currentFloor	 	int

func setLights(buttonLights [src.N_FLOORS][3]int, chCommandTI chan io.Command){
	for floor := 0; floor < src.N_FLOORS; floor ++ {
		chCommandTI <- io.Command{io.SET_BUTTON_LAMP,buttonLights[floor][src.BUTTON_UP],     floor,src.BUTTON_UP}
		chCommandTI <- io.Command{io.SET_BUTTON_LAMP,buttonLights[floor][src.BUTTON_DOWN],   floor,src.BUTTON_DOWN}
		chCommandTI <- io.Command{io.SET_BUTTON_LAMP,buttonLights[floor][src.BUTTON_INSIDE], floor,src.BUTTON_INSIDE}
	}
}

func goUpUntilReachFloor(chCommandTI chan io.Command, chFloorFI chan int){
	chCommandTI  	<- io.Command{io.SET_MOTOR_DIR,src.DIR_UP,-1,src.BUTTON_NONE}
	currentFloor	= <- chFloorFI
	chCommandTI  	<- io.Command{io.SET_MOTOR_DIR,src.DIR_STOP,-1,src.BUTTON_NONE}
	chCommandTI  	<- io.Command{io.SET_FLOOR_INDICATOR_LAMP,src.ON,currentFloor,src.BUTTON_NONE}
}

func doorTimer(chOpenDoor chan bool, chCloseDoor chan bool){
	
	for{
		<-chOpenDoor
		time.Sleep(2*time.Second)
		chCloseDoor <- true		
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

func InitController() {
	
	chCommandTI 			:= make(chan io.Command)
	chOrderFI 				:= make(chan src.ButtonOrder,512)
	chFloorFI				:= make(chan int)

	chFloorTQ				:= make(chan int,2)
	chOrderTQ 				:= make(chan src.ButtonOrder, 512)
	chOrderFinishedTQ 		:= make(chan bool,2)
	chButtonLightsFQ 		:= make(chan [src.N_FLOORS][3]int,2)
	chDestinationFloorFQ	:= make(chan int,2)
	
	chOpenDoor 				:= make(chan bool)
	chCloseDoor				:= make(chan bool)

	io.InitIo(chCommandTI, chOrderFI, chFloorFI)
	goUpUntilReachFloor(chCommandTI, chFloorFI)	
	
	state = IDLE
	
	go queue.InitQueue(chFloorTQ, chOrderTQ, chOrderFinishedTQ, chButtonLightsFQ, chDestinationFloorFQ)
	chFloorTQ <- currentFloor

	go controllerManager(chCommandTI, chOrderFI, chFloorFI, chFloorTQ, chOrderTQ, chOrderFinishedTQ, chButtonLightsFQ, chDestinationFloorFQ,chOpenDoor,chCloseDoor)
	go doorTimer(chOpenDoor, chCloseDoor)
}

func controllerManager( chCommandTI chan io.Command,
					 	chOrderFI chan src.ButtonOrder,
					 	chFloorFI chan int, 
					 	chFloorTQ chan int, 
					 	chOrderTQ chan src.ButtonOrder,
					 	chOrderFinishedTQ chan bool, 
					 	chButtonLightsFQ chan [src.N_FLOORS][3]int,
					 	chDestinationFloorFQ chan int, 
					 	chOpenDoor chan bool, 
					 	chCloseDoor chan bool){
	for {
		select {
			
			case order := <- chOrderFI:
				 chOrderTQ <- order

			case currentFloor = <- chFloorFI:
				 chCommandTI  	<- io.Command{io.SET_FLOOR_INDICATOR_LAMP,src.ON, currentFloor, src.BUTTON_NONE}
				 chFloorTQ 		<- currentFloor
				 destinationFloor = <- chDestinationFloorFQ
				
				switch state{
					case MOVING:
						elevatorDirection   := chooseElevatorDirection()
						chCommandTI 		<- io.Command{io.SET_MOTOR_DIR,elevatorDirection,-1,src.BUTTON_NONE}

						if(currentFloor == destinationFloor){
							state = DOOR_OPEN
							chCommandTI <- io.Command{io.SET_DOOR_OPEN_LAMP,src.ON,-1,src.BUTTON_NONE}
							chOpenDoor  <- true
						}

					default:
						continue
				}
				
				
			case <- chCloseDoor:
					state = IDLE
					chCommandTI <- io.Command{io.SET_DOOR_OPEN_LAMP,src.OFF,-1,src.BUTTON_NONE}
					chOrderFinishedTQ <- true

			case buttonLights := <-chButtonLightsFQ:
				setLights(buttonLights,chCommandTI)
		
			case destinationFloor = <- chDestinationFloorFQ:
				switch state{
					case IDLE:
						elevatorDirection := chooseElevatorDirection()
						
						if(currentFloor == destinationFloor){
							state = DOOR_OPEN
							chCommandTI <- io.Command{io.SET_DOOR_OPEN_LAMP,src.ON,-1,src.BUTTON_NONE}
							chOpenDoor <- true
						
						}else{
							chCommandTI <- io.Command{io.SET_MOTOR_DIR,elevatorDirection,-1,src.BUTTON_NONE}
							state = MOVING
						}
					default:
						continue
				}		
		}
	}
}