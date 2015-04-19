package queue

import (
	//"strings"
	//"strconv"
	"Heisprosjekt/src"
	"Heisprosjekt/network"
	"time"
	"Heisprosjekt/tools"
)

const (
	ORDER = 1
	EMPTY = 0
	DELETE_ORDER = -1
)

var ourID			string
var globalOrders    src.ElevatorData
var currentOrder	src.ButtonOrder
var elevatorQueues	= make(map[string] src.ElevatorData)
var timeoutLimit time.Duration = 1*time.Second


// TESTED:
func mergeOrders(elevatorQueues map[string] src.ElevatorData) src.ElevatorData {
	var mergedData src.ElevatorData

	for floor := 0; floor < src.N_FLOORS; floor ++ {
		for direction := 0; direction < 2; direction ++ {
			for _, elevatorQueue := range elevatorQueues {
				if (elevatorQueue.OutsideOrders[floor][direction] == src.ORDER) {
					mergedData.OutsideOrders[floor][direction] = src.ORDER
				}
			}
		}
	}

	mergedData.InsideOrders = elevatorQueues[ourID].InsideOrders
	return mergedData
}

// TESTED:
func assignOrder(elevatorQueues map[string] src.ElevatorData, order src.ButtonOrder) {
	var minID string
	minCost := 10000 // skulle gjerne vært inf elno..

	for elevatorID, elevatorQueue := range elevatorQueues {
		cost := calcOrderCost(elevatorQueue, order)
		if (cost < minCost) {
			minID = elevatorID
			minCost = cost
		
		} else if (cost == minCost && elevatorID < minID) {
			minID = elevatorID
			minCost = cost
		}
	}

	if (minID == ourID) {
		addOrder(ourID, order)
	}
}


func calcOrderCost(elevator src.ElevatorData, order src.ButtonOrder) int {
	return abs(elevator.Floor - order.Floor)
}

// TO BE EDITED!
func calcNextOrderAndFloor(queueMatrix src.ElevatorData) src.ButtonOrder {
	currentOrder := src.ButtonOrder{-1, src.BUTTON_NONE}

	switch queueMatrix.Direction {
		case src.DIR_UP, src.DIR_STOP:
			
			for floor := queueMatrix.Floor; floor < src.N_FLOORS; floor ++ {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_UP] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_UP}
					return currentOrder
				}
			}

			for floor := src.N_FLOORS - 1; floor >= 0; floor -- {
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
			
			for floor := queueMatrix.Floor; floor >= 0; floor -- {
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


func addOrder(ID string, order src.ButtonOrder) {
	switch(order.ButtonType) {

		case src.BUTTON_INSIDE:
			tmp := elevatorQueues[ID]
			tmp.InsideOrders[order.Floor] = ORDER
			elevatorQueues[ID] = tmp
		
		default:
			tmp := elevatorQueues[ID]
			tmp.OutsideOrders[order.Floor][order.ButtonType] = ORDER
			elevatorQueues[ID] = tmp
	}
}


func deleteOrder(ID string, order src.ButtonOrder) {
	if (ButtonOrder.Floor != -1 && ButtonOrder.ButtonType != src.BUTTON_NONE) {
		tmp := elevatorQueues[ID]
		tmp.InsideOrders[order.Floor] = EMPTY
		tmp.OutsideOrders[order.Floor][order.ButtonType] = EMPTY
		elevatorQueues[ID] = tmp
	}
}


// func InitQueue(chFloorFromController chan int, chOrderFromController chan src.ButtonOrder, chDirectionFromController chan int, chFinishedFromController chan bool, chGlobalOrdersToController chan src.ElevatorData, chDestinationFloorToController chan int) {
// 	go network.NetworkHandler()
// 	ourID = <- network.ChIDFromNetwork
// 	var InitialQueue src.ElevatorData

// 	InitialQueue.Floor = <- chFloorFromController
// 	println("lol")
// 	elevatorQueues[ourID] = InitialQueue
// 	network.ChQueueReadyToBeSent <- elevatorQueues[ourID]

// 	go queueManager(chFloorFromController, 
// 					chOrderFromController, 
// 					chDirectionFromController, 
// 					chFinishedFromController, 
// 					chGlobalOrdersToController, 
// 					chDestinationFloorToController)
// }

func timer(timeout chan bool) {
	for {
		time.Sleep(timeoutLimit)
		timeout <- true
	}
}

func QueueManager(chFloorFromController chan int, chOrderFromController chan src.ButtonOrder, chDirectionFromController chan int, chFinishedFromController chan bool, chGlobalOrdersToController chan src.ElevatorData, chDestinationFloorToController chan int) {
	var chUpdateGlobalOrders = make(chan bool)

	go network.NetworkHandler()
	ourID = <- network.ChIDFromNetwork
	var InitialQueue src.ElevatorData

	InitialQueue.Floor = <- chFloorFromController
	elevatorQueues[ourID] = InitialQueue

	network.ChQueueReadyToBeSent <- elevatorQueues[ourID]	
	go timer(chUpdateGlobalOrders)

	for {
		select {
			
			case <- chUpdateGlobalOrders:
				globalOrders = mergeOrders(elevatorQueues)
				chGlobalOrdersToController <- globalOrders
				println(" ---- ORDERS ------")
				tools.PrintQueue(elevatorQueues[ourID])

			case order := <- network.ChOrderToQueue:
				assignOrder(elevatorQueues, order)
				currentOrder = calcNextOrderAndFloor(elevatorQueues[ourID])
				if currentOrder.Floor != -1 {chDestinationFloorToController <- currentOrder.Floor}
				network.ChQueueReadyToBeSent <- elevatorQueues[ourID]

			case updatedQueue := <- network.ChElevatorDataToQueue:
				elevatorQueues[updatedQueue.SenderAddress] = updatedQueue.Data

			case floor := <- chFloorFromController:
				tmp := elevatorQueues[ourID]
				tmp.Floor = floor
				elevatorQueues[ourID] = tmp
				network.ChQueueReadyToBeSent <- elevatorQueues[ourID]

			case order := <- chOrderFromController:
				network.ChOrderFromQueue <- order
				assignOrder(elevatorQueues, order)
				currentOrder = calcNextOrderAndFloor(elevatorQueues[ourID])
				println("sending nextfloor to C")
				if currentOrder.Floor != -1 {chDestinationFloorToController <- currentOrder.Floor}
				println("... sent!")
				network.ChQueueReadyToBeSent <- elevatorQueues[ourID]

			case <- chFinishedFromController:
				deleteOrder(ourID, currentOrder)
				network.ChQueueReadyToBeSent <- elevatorQueues[ourID]
				currentOrder = calcNextOrderAndFloor(elevatorQueues[ourID])
				if currentOrder.Floor != -1 {chDestinationFloorToController <- currentOrder.Floor}

			case direction := <- chDirectionFromController:
				tmp := elevatorQueues[ourID]
				tmp.Direction = direction
				elevatorQueues[ourID] = tmp
				network.ChQueueReadyToBeSent <- elevatorQueues[ourID]

			case elevator := <- network.ChLostElevator:
				dataToDistrubute := elevatorQueues[elevator]
				delete(elevatorQueues, elevator)
				
				for floor := 0; floor < src.N_FLOORS; floor ++ {
					for direction := 0; direction < 2; direction ++ {
						if (dataToDistrubute.OutsideOrders[floor][direction] == src.ORDER) {
							assignOrder(elevatorQueues, src.ButtonOrder{floor, direction})
						}
					}
				}

				network.ChQueueReadyToBeSent <- elevatorQueues[ourID]
				currentOrder = calcNextOrderAndFloor(elevatorQueues[ourID])
				if currentOrder.Floor != -1 {chDestinationFloorToController <- currentOrder.Floor}

			default:
				time.Sleep(100*time.Millisecond)
		}
	}
}

func abs(value int) int {
	if (value >= 0) { 	return value
	}else { 				return -1*value}
}

// TODO:
// * renew calcNextFloor and add functionality for DIR_STOP												|

// BUGS:
// * potensielt: kan to heiser gå til samme etasje, en med inn og en med opp? I så fall blir det krøll.	|
// * Hvis vi står i 2. etg, og behandler "opp" ordre, og dermed trykker på "opp"-ordre i 1. etg, vil 
//   heisen slette ordren i 1... Dette er fordi fnishedOrder blir sendt, og currentOrder er da i 1.		|