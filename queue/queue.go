package queue

import (
	//"strings"
	//"strconv"
	"Heisprosjekt/src"
	"Heisprosjekt/network"
	"time"
	//"Heisprosjekt/tools"
)

const (
	ORDER = 1
	EMPTY = 0
	DELETE_ORDER = -1
)

var ourID			string
var globalOrders    src.ElevatorData
var currentOrder	src.ButtonOrder
var elevatorQueues	= make(map[string] src.ElevatorData) // Must use pointer to let golang be able to edit the map! :(

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
	tmp := elevatorQueues[ID]
	tmp.InsideOrders[order.Floor] = EMPTY
	tmp.OutsideOrders[order.Floor][order.ButtonType] = EMPTY
	elevatorQueues[ID] = tmp
}


func InitQueue(ChFloorFromController chan int, chOrderFromController chan src.ButtonOrder, chDirectionFromController chan int, chFinishedFromController chan bool, chGlobalOrdersToController chan src.ElevatorData, chNewNextFloorFromQueue chan int) {
	//go network.NetworkHandler()
	//ourID := network.GetIPAddress()
	var emptyQueue src.ElevatorData
	ourID = "192.168.0.102"
	elevatorQueues[ourID] = emptyQueue

	go queueManager(ChFloorFromController, 
					chOrderFromController, 
					chDirectionFromController, 
					chFinishedFromController, 
					chGlobalOrdersToController, 
					chNewNextFloorFromQueue)
}


func queueManager(ChFloorFromController chan int, chOrderFromController chan src.ButtonOrder, chDirectionFromController chan int, chFinishedFromController chan bool, chGlobalOrdersToController chan src.ElevatorData, chNewNextFloorFromQueue chan int) {
	for {
		select {
			
			case <- network.ChReadyToMerge:
				globalOrders = mergeOrders(elevatorQueues)
				chGlobalOrdersToController <- globalOrders

			case order := <- network.ChOrderToQueue:
				assignOrder(elevatorQueues, order)
				currentOrder = calcNextOrderAndFloor(elevatorQueues[ourID])
				chNewNextFloorFromQueue <- currentOrder.Floor

			case updatedQueue := <- network.ChQueueMessage:
				elevatorQueues[updatedQueue.SenderAddress] = updatedQueue.Data

			case floor := <- ChFloorFromController:
				tmp := elevatorQueues[ourID]
				tmp.Floor = floor
				elevatorQueues[ourID] = tmp

			case order := <- chOrderFromController:
				//network.ChOrderFromQueue <- order
				assignOrder(elevatorQueues, order)
				currentOrder = calcNextOrderAndFloor(elevatorQueues[ourID])
				chNewNextFloorFromQueue <- currentOrder.Floor


			case <- chFinishedFromController:
				deleteOrder(ourID, currentOrder)

			case direction := <- chDirectionFromController:
				tmp := elevatorQueues[ourID]
				tmp.Direction = direction
				elevatorQueues[ourID] = tmp

			// case elevator has disconnected!

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

// * 
// * pass our queue into merge orders when ReadyToMerge 					| OK
// * idea: make cost-function to weight no. of stops only (less code)		|

// BUGS:

// * potensielt: kan to heiser gå til samme etasje, en med inn og en med opp? I så fall blir det krøll.	|
// * Hvis vi står i 2. etg, og behandler "opp" ordre, og dermed trykker på "opp"-ordre i 1. etg, vil 
//   heisen slette ordren i 1... Dette er fordi fnishedOrder blir sendt, og currentOrder er da i 1.		|
//


// TODO:
// * renew calcNextFloor and add functionality for DIR_STOP					|

