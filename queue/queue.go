package queue

import (
	"strings"
	"strconv"
	"Heisprosjekt/src"
	"Heisprosjekt/network"
)

const (
	ORDER = 1
	EMPTY = 0
	DELETE_ORDER = -1
)

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


func clearOutsideOrders(queueList []src.ElevatorData) {
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


func assignOrders(queueList []src.ElevatorData, mergedQueue src.ElevatorData) src.ElevatorData {
	var costArray []int
	
	clearOutsideOrders(queueList)
	
	for direction := 0; direction < 2; direction ++ {
		for floor := 0; floor < src.N_FLOORS; floor ++ {
				
			if mergedQueue.OutsideOrders[floor][direction] == src.ORDER {

				for eachElevator := 0; eachElevator < len(queueList); eachElevator ++{
					
					queueList[eachElevator].OutsideOrders[floor][direction] = ORDER
					costArray = append(costArray, calcTotalCost(&queueList[eachElevator]))
					queueList[eachElevator].OutsideOrders[floor][direction] = EMPTY
				}

				assignedElevator := findMinimumCost(costArray, queueList)
				queueList[assignedElevator].OutsideOrders[floor][direction] = ORDER
				costArray = []int{}
			}
		}
	}

	ourQueue := queueList[len(queueList) - 1]
	return ourQueue
}

// TESTED:
func calcCostUp(queueData src.ElevatorData, COST_FOR_MOVE int, COST_FOR_ORDER int) int {
	totalCost := 0
	floorsSinceLastOrder := 0

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
}

// TESTED:
func calcCostDown(queueData src.ElevatorData, COST_FOR_MOVE int, COST_FOR_ORDER int) int {
	totalCost := 0
	floorsSinceLastOrder := 0
	
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
}

// TESTED:
func calcTotalCost(queueData *src.ElevatorData) int {
	const COST_FOR_ORDER = 3
	const COST_FOR_MOVE = 1

	switch (*queueData).Direction {
		case src.DIR_UP:
			return calcCostUp(*queueData, COST_FOR_MOVE, COST_FOR_ORDER)
			
		case src.DIR_DOWN:
			return calcCostDown(*queueData, COST_FOR_MOVE, COST_FOR_ORDER)
			
		default:
			// Finnie min av costup og costdown og sette dir
			costUp := calcCostUp(*queueData, COST_FOR_MOVE, COST_FOR_ORDER)
			costDown := calcCostDown(*queueData, COST_FOR_MOVE, COST_FOR_ORDER)

			if costUp == 0 && costDown == 0 {
				return 0
			}

			if costUp > costDown {
				(*queueData).Direction = src.DIR_DOWN
				return costDown
			}
			
			if costUp <= costDown {
				(*queueData).Direction = src.DIR_UP
				return costUp
			}

			// En liten bug: Hvis vi står en topp- eller bunn-etasje, vil costUp = costDown.
			// Test og se om det kan være sånn, eller om det må fikses.
			return 0;
	}
}


func calcNextFloor(queueMatrix src.ElevatorData) int {
	nextFloor := queueMatrix.Floor

	switch queueMatrix.Direction {
		case src.DIR_UP:
			
			for floor := queueMatrix.Floor; floor < src.N_FLOORS; floor ++ {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_UP] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					nextFloor = floor
					return nextFloor
				}
			}

			for floor := src.N_FLOORS - 1; floor > 0; floor -- {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_DOWN] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					nextFloor = floor
					return nextFloor
				}
			}

			for floor := 0; floor < queueMatrix.Floor; floor ++ {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_DOWN] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					nextFloor = floor
					return nextFloor
				}
			}
		case src.DIR_DOWN:
			
			for floor := queueMatrix.Floor; floor > 0; floor -- {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_DOWN] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					nextFloor = floor
					return nextFloor
				}
			}

			for floor := 0; floor < src.N_FLOORS; floor ++ {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_UP] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					nextFloor = floor
					return nextFloor
				}
			}

			for floor := src.N_FLOORS - 1; floor > queueMatrix.Floor; floor -- {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_DOWN] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					nextFloor = floor
					return nextFloor
				}
			}

		default:
			return queueMatrix.Floor
	}
	return nextFloor
}


func QueueHandler() {
	chNewFloor		:= make(chan int)
	chNewOrder 		:= make(chan src.ButtonOrder)
	chNewDirection	:= make(chan int)
	chOrderIsFinished 	:= make(chan src.ButtonOrder)

	IPaddr := network.GetIPAddress()
	IPaddrArray := strings.Split(IPaddr, ".")
	knownOrders.ID, _ = strconv.Atoi(IPaddrArray[3])

	for {
		select {
			case <- network.ChReadyToMerge:
				allElevatorData := append(listOfIncomingData, knownOrders)
				mergedQueue := mergeOrders(allElevatorData)
				knownOrders.OutsideOrders = mergedQueue.OutsideOrders
				network.ChQueueReadyToBeSent <- knownOrders
				assignedOrder := assignOrders(allElevatorData, mergedQueue)
				calcNextFloor(assignedOrder)
				listOfIncomingData = nil

			case data := <- network.ChDataToQueue:
				listOfIncomingData = append(listOfIncomingData, data)

			case newFloor := <- chNewFloor:
				knownOrders.Floor = newFloor

			case newOrder := <- chNewOrder:
				//update elevatorData.queueMatrix
				print(newOrder.Floor) // ONLY FOR TESTING
			case finishedOrder := <- chOrderIsFinished:
				//update elevatorData.queueMatrix with a special char
				print(finishedOrder.Floor) // ONLY FOR TESTING
			case newDirection := <- chNewDirection:
				knownOrders.Direction = newDirection

		}
	}
}

// TODO:

// * implement functionality for DIR_STOP in calcTotalCost					| OK
// * pass our queue into merge orders when ReadyToMerge 					| OK
// * idea: make cost-function to weight no. of stops only (less code)		|
// * implement addOrder() and deleteOrder()									|