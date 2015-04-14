package queue

import (
	"strings"
	"strconv"
	"Heisprosjekt/src"
	"Heisprosjekt/network"
	"Heisprosjekt/tools"
)

const (
	ORDER = 1
	EMPTY = 0
	DELETE_ORDER = -1
)

var knownOrders    		src.ElevatorData
var listOfIncomingData  []src.ElevatorData
var currentOrder		src.ButtonOrder

func findMinimumCost(costArray []int, queueList []src.ElevatorData) int {
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

	for floor := 0; floor < src.N_FLOORS; floor ++ {
		for direction := 0; direction < 2; direction ++ {

			nextDirection:
			for eachQueue := 0; eachQueue < len(queueList); eachQueue ++ {
				switch queueList[eachQueue].OutsideOrders[floor][direction] {
					case src.ORDER:
						mergedData.OutsideOrders[floor][direction] = src.ORDER
					case src.DELETE_ORDER:
						mergedData.OutsideOrders[floor][direction] = src.EMPTY
						break nextDirection
				}
			}
		}
	}

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


func calcNextOrderAndFloor(queueMatrix src.ElevatorData) src.ButtonOrder {
	currentOrder := src.ButtonOrder{-1, src.BUTTON_NONE}

	switch queueMatrix.Direction {
		case src.DIR_UP:
			
			for floor := queueMatrix.Floor; floor < src.N_FLOORS; floor ++ {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_UP] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_UP}
					return currentOrder
				}
			}

			for floor := src.N_FLOORS - 1; floor > 0; floor -- {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_DOWN] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_DOWN}
					return currentOrder
				}
			}

			for floor := 0; floor < queueMatrix.Floor; floor ++ {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_UP] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_UP}
					return currentOrder
				}
			}
		case src.DIR_DOWN:
			
			for floor := queueMatrix.Floor; floor > 0; floor -- {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_DOWN] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_DOWN}
					return currentOrder
				}
			}

			for floor := 0; floor < src.N_FLOORS; floor ++ {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_UP] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_UP}
					return currentOrder
				}
			}

			for floor := src.N_FLOORS - 1; floor > queueMatrix.Floor; floor -- {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_DOWN] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_DOWN}
					return currentOrder
				}
			}

		default:
			return currentOrder
	}
	return currentOrder
}


func RememberDeletedOrders(queueData src.ElevatorData) (src.ElevatorData, int) {
	var memoOfDeletedOrders src.ElevatorData
	storedDeletedOrder := 0
	
	for floor := 0; floor < src.N_FLOORS; floor ++ {
		for direction := 0; direction < 2; direction ++ {
			if queueData.OutsideOrders[floor][direction] == DELETE_ORDER {
				memoOfDeletedOrders.OutsideOrders[floor][direction] = DELETE_ORDER
				storedDeletedOrder = 1
			}
		}

		if queueData.InsideOrders[floor] == DELETE_ORDER {
			memoOfDeletedOrders.InsideOrders[floor] = DELETE_ORDER
			storedDeletedOrder = 1
		}
	}

	return memoOfDeletedOrders, storedDeletedOrder
}


func addDeletedOrders(queueData src.ElevatorData, memoOfDeletedOrders src.ElevatorData) src.ElevatorData {
	for floor := 0; floor < src.N_FLOORS; floor ++ {
		for direction := 0; direction < 2; direction ++ {
			if memoOfDeletedOrders.OutsideOrders[floor][direction] == DELETE_ORDER {
				queueData.OutsideOrders[floor][direction] = DELETE_ORDER
			}
		}

		if memoOfDeletedOrders.InsideOrders[floor] == DELETE_ORDER {
			queueData.InsideOrders[floor] = DELETE_ORDER
		}
	}

	return queueData
}


func InitQueue(chNewFloor chan int, chNewOrder chan src.ButtonOrder, chNewDirection chan int, chOrderIsFinished chan bool, chNewOrdersFromQueue chan src.ElevatorData, chNewNextFloorFromQueue chan int) {
	go network.NetworkHandler()

	IPaddr := network.GetIPAddress()
	IPaddrArray := strings.Split(IPaddr, ".")
	knownOrders.ID, _ = strconv.Atoi(IPaddrArray[3])

	storedDeletedOrder := 0
	var memoOfDeletedOrders src.ElevatorData

	go QueueHandler(chNewFloor, chNewOrder, chNewDirection, chOrderIsFinished, chNewOrdersFromQueue, chNewNextFloorFromQueue, storedDeletedOrder, memoOfDeletedOrders)
}


func addOrder(order src.ButtonOrder) {
	switch(order.ButtonType) {
		case src.BUTTON_INSIDE:
			knownOrders.InsideOrders[order.Floor] = ORDER
		default:
			knownOrders.OutsideOrders[order.Floor][order.ButtonType] = ORDER


	}
}


func deleteOrder(order src.ButtonOrder) {
	knownOrders.InsideOrders[order.Floor] = EMPTY
	knownOrders.OutsideOrders[order.Floor][order.ButtonType] = DELETE_ORDER
}


func QueueHandler(chNewFloor chan int, chNewOrder chan src.ButtonOrder, chNewDirection chan int, chOrderIsFinished chan bool, chNewOrdersFromQueue chan src.ElevatorData, chNewNextFloorFromQueue chan int, storedDeletedOrder int, memoOfDeletedOrders src.ElevatorData) {
	
	for {
		select {
			case <- network.ChReadyToMerge:
				allElevatorData := append(listOfIncomingData, knownOrders)
				
				//println("------------------")
				// println("Registered Queues:")
				// tools.PrintQueueArray(allElevatorData)

				// THis fixes the bug with deleted order may not be registered correctly.
				if storedDeletedOrder == 0 { 
					memoOfDeletedOrders, storedDeletedOrder = RememberDeletedOrders(knownOrders)
				} else if storedDeletedOrder == 1 { 
					storedDeletedOrder ++ 
				} else if storedDeletedOrder == 2 {
					storedDeletedOrder = 0
				} else {}
				// --------------------

				network.ChQueueReadyToBeSent <- knownOrders
				mergedQueue := mergeOrders(allElevatorData)
				knownOrders.OutsideOrders = mergedQueue.OutsideOrders
				assignedOrder := assignOrders(allElevatorData, mergedQueue)
				currentOrder = calcNextOrderAndFloor(assignedOrder)
				listOfIncomingData = nil

				if (currentOrder.Floor != -1) { chNewNextFloorFromQueue <- currentOrder.Floor}
				chNewOrdersFromQueue <- knownOrders
				if storedDeletedOrder > 0 { knownOrders = addDeletedOrders(knownOrders, memoOfDeletedOrders)}
				
				//tools.PrintQueueHandler(mergedQueue, assignedOrder)
				tools.PrintQueue(assignedOrder)
				println("Next Floor:", currentOrder.Floor)

			case data := <- network.ChDataToQueue:
				listOfIncomingData = append(listOfIncomingData, data)

			case newFloor := <- chNewFloor:
				knownOrders.Floor = newFloor

			case newOrder := <- chNewOrder:
				addOrder(newOrder)

			case <- chOrderIsFinished:
				deleteOrder(currentOrder)
				println("Order is deleted!")

			case newDirection := <- chNewDirection:
				knownOrders.Direction = newDirection
		}
	}
}


// TODO:

// * implement functionality for DIR_STOP in calcTotalCost					| OK
// * pass our queue into merge orders when ReadyToMerge 					| OK
// * idea: make cost-function to weight no. of stops only (less code)		|
// * implement addOrder() and deleteOrder()									| OK
// * remove knownOrders and incomingList as global var						|
// * make addDeletedOrders be able to remember multiple deletions			|

// BUGS:

// * potensielt: kan to heiser gå til samme etasje, en med inn og en med opp? I så fall blir det krøll.
