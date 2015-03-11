

package main//package queue

type elevatorData struct { 
	floor 			int
	direction 		int
	queueMatrix  [][]int
}

type Order struct {	// THis struct should be imported from styring / IO!
	floor 			int
	buttonType		int
}	


//Definerer egne variabler:
var localQueue     	elevatorData
var globalQueues  []elevatorData
var FLOORS 		   	int //Import from Styring!


func mergeOrders(queueList []elevatorData) elevatorData {
	//Go through list, append orders. Remove if orders are finished
	mergedQueue := [FLOORS][2]int	
}

func assignOrders(queueList []elevatorData, mergedQueue elevatorData) elevatorData {
	//queueList = globalQueues + localQueue
	//for direction = 0; direction < 2; direction ++:
		//for floor = 0; floor < FLOOR; floor ++:
			//if mergedQueue.queueMatrix[floor][direction] == 1:
				//Altsaa, hvis bestilling ...
				//for elevator = 0; elevator < len(queueList); elevator ++:
					//calcTotalCost(elevator)

}

func calcTotalCost() int {

}

func calcNextFloor() {

}


func main() { // func queueHandler() {
	// Initialisere kanaler her:
	chReadyToMerge 	:= make(chan bool)
	chIncomingQueues:= make(chan elevatorData)
	chNewFloor		:= make(chan int)
	chNewOrder 		:= make(chan Order) //make(chan styring.Order)
	chNewDirection	:= make(chan int)

	for {
		select {
			case <- chReadyToMerge:
				//add our elevatorData into dataStorage
				//mergeOrders()
				//assignOrders()
				//calcNextFloor()
				//send mergedOrders til nettverk
				//toem dataStorage
			case data := <- chIncomingQueues:
				//append to list of data
			case newFloor := <- chNewFloor:
				//update elevatorData.floor
			case newOrder := <- chNewOrder:
				//update elevatorData.queueMatrix
			case deletedOrder := <- chDeleteOrder:
				//update elevatorData.queueMatrix with a special char
			case newDirection := <- chNewDirection:
				//update elevator.direction

		}
	}
}
