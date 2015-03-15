package queue

import (
	"Heisprosjekt/src"
	"Heisprosjekt/network"
)

const (
	ORDER = 1
	EMPTY = 0
	DELETE_ORDER = -1
)

var ID 					int // The IP has to be converted into int somehow
var knownOrders    		src.ElevatorData
var listOfIncomingData  []src.ElevatorData


func findMinimumCost(costArray []int, queueList []src.ElevatorData) int {
	// This function is now influenced by a random factor - when two costs
	// are equal, the first entry is set as minimum. This cannot be the case, 
	// as two elevators may calculate different!
	minValue := costArray[0]
	minElevator := 0

	for elevator, value := range costArray {
		if value < minValue {
			minValue = value
			minElevator = elevator
		}

		if value == minValue {
			if queueList[elevator].ID < queueList[minElevator].ID {
				minValue = value
				minElevator = elevator
			}
		}
	}
	return minElevator
}


func clearOutsideOrders(queueList []src.ElevatorData, ourQueue src.ElevatorData) {
	var EmptyElevatorData src.ElevatorData

	for elevator := 0; elevator < len(queueList); elevator ++ {
		queueList[elevator].OutsideOrders = EmptyElevatorData.OutsideOrders
	}
}

// TESTED:
func mergeOrders(queueList []src.ElevatorData) src.ElevatorData {
	var mergedData src.ElevatorData
	var mergedQueue [src.N_FLOORS][2]int	

	for floor := 0; floor < src.N_FLOORS; floor ++ {
		directionLoop:	
		for direction := 0; direction < 2; direction ++ {
			
			for eachQueue := 0; eachQueue < len(queueList); eachQueue ++ {
				switch queueList[eachQueue].OutsideOrders[floor][direction] {
					case src.ORDER:
						mergedQueue[floor][direction] = src.ORDER
					case src.DELETE_ORDER:
						mergedQueue[floor][direction] = src.EMPTY
						break directionLoop
				}
			}
		}
	}

	mergedData.OutsideOrders = mergedQueue
	return mergedData
}


func assignOrders(queueList []src.ElevatorData, mergedQueue src.ElevatorData) []src.ElevatorData {
	var costArray []int
	ourQueue := knownOrders

	clearOutsideOrders(queueList, ourQueue)
	
	for direction := 0; direction < 2; direction ++ {
		for floor := 0; floor < src.N_FLOORS; floor ++ {
				
			if mergedQueue.OutsideOrders[floor][direction] == src.ORDER {

				for eachElevator := 0; eachElevator < len(queueList); eachElevator ++{
					
					queueList[eachElevator].OutsideOrders[floor][direction] = ORDER
					costArray = append(costArray, calcTotalCost(queueList[eachElevator]))
					queueList[eachElevator].OutsideOrders[floor][direction] = EMPTY
				}

				assignedElevator := findMinimumCost(costArray, queueList)
				queueList[assignedElevator].OutsideOrders[floor][direction] = ORDER
				costArray = []int{}
			}
		}
	}
	// Loop through and search for our list, alternatively we know that our is the last one.
	// Then return ourQueue only
	return queueList
}

// Need to implement when direction is DIR_STOP
func calcTotalCost(queueData src.ElevatorData) int {
	const COST_FOR_ORDER = 3
	const COST_FOR_MOVE = 1
	totalCost := 0
	floorsSinceLastOrder := 0

	// IDE:
	// Kan man gjoere det saa enkelt at man teller oppover og nedover helt til man finner siste 1-er
	// Deretter teller antall 1-ere?
	switch queueData.Direction {
		case src.DIR_UP:

			for floor := queueData.Floor; floor < src.N_FLOORS; floor ++ {
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1
				if queueData.OutsideOrders[floor][0] == src.ORDER || queueData.InsideOrders[floor] == src.ORDER {
					totalCost += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}				
			}

			floorsSinceLastOrder -= 1 

			for floor := src.N_FLOORS-1; floor >= 0; floor-- {
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1 
				if queueData.OutsideOrders[floor][1] == src.ORDER || (queueData.InsideOrders[floor] == src.ORDER  && floor < queueData.Floor){
					totalCost += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}				
			}
			
			floorsSinceLastOrder -= 1

			for floor := 0; floor < queueData.Floor; floor ++{
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1
				if queueData.OutsideOrders[floor][0] == src.ORDER {
					totalCost  += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}
			}

			totalCost -= 3*COST_FOR_MOVE // 1 when we start, and 1 when we change dir * 2
			totalCost -= floorsSinceLastOrder
			return totalCost
		
		case src.DIR_DOWN:

			for floor := queueData.Floor; floor >= 0; floor-- {
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1 
				if queueData.OutsideOrders[floor][1] == src.ORDER || queueData.InsideOrders[floor] == src.ORDER  {
					totalCost += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}				
			}
 
			floorsSinceLastOrder -= 1 

			for floor := 0; floor < src.N_FLOORS; floor ++ {
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1
				if queueData.OutsideOrders[floor][0] == src.ORDER || (queueData.InsideOrders[floor] == src.ORDER && floor > queueData.Floor){
					totalCost += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}				
			}

			floorsSinceLastOrder -= 1

			for floor := src.N_FLOORS-1; floor > queueData.Floor; floor --{
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1
				if queueData.OutsideOrders[floor][1] == src.ORDER {
					totalCost  += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}
			}

			totalCost -= 3*COST_FOR_MOVE // 1 when we start, and 1 when we change dir * 2
			totalCost -= floorsSinceLastOrder
			return totalCost

		default: //Ta med src.DIR_STOP?
			// ...
	}

	// Gaa gjennom hele koen (bruk dir for aa bestemme start-retning) og legg til costs
	// for hver etasje/ordre helt til alle ordrene er tatt med!
	return 0;
}


func calcNextFloor() {
}


func QueueHandler() {
	// Initialisere kanaler her:
	chNewFloor		:= make(chan int)
	chNewOrder 		:= make(chan src.ButtonOrder) //make(chan styring.Order)
	chNewDirection	:= make(chan int)
	chOrderIsFinished 	:= make(chan src.ButtonOrder) //make(chan styring.Order)

	for {
		select {
			case <- network.ChReadyToMerge:
				mergedQueue := mergeOrders(listOfIncomingData)
				listOfIncomingData = nil
				network.ChQueueReadyToBeSent <- mergedQueue
				//assignOrders()
				//calcNextFloor()
				//toem dataStorage
			case data := <- network.ChDataToQueue:
				listOfIncomingData = append(listOfIncomingData, data)

			case newFloor := <- chNewFloor:
				//update elevatorData.floor
				print(newFloor) // ONLY FOR TESTING
			case newOrder := <- chNewOrder:
				//update elevatorData.queueMatrix
				print(newOrder.Floor) // ONLY FOR TESTING
			case finishedOrder := <- chOrderIsFinished:
				//update elevatorData.queueMatrix with a special char
				print(finishedOrder.Floor) // ONLY FOR TESTING
			case newDirection := <- chNewDirection:
				//update elevator.direction
				print(newDirection) // ONLY FOR TESTING

		}
	}
}

// TODO:

// * implement functionality for DIR_STOP in calcTotalCost					|
// * pass our queue into merge orders when ReadyToMerge 					|
// * idea: make cost-function to weight no. of stops only (less code)		|