package queue

import (
	"Heisprosjekt/src"
	"Heisprosjekt/network"
	"time"
	"Heisprosjekt/tools"
)

const (
	ORDER = 1
	NO_ORDER = 0
	COST_FOR_ORDER = 2
	COST_FOR_MOVEMENT = 1
)

var ourID			string
var buttonLights    src.ElevatorData
var currentOrder	src.ButtonOrder
var elevatorQueues 	= make(map[string] src.ElevatorData)
var timeoutLimit time.Duration = 100*time.Millisecond

func determineButtonLights(elevatorQueues map[string] src.ElevatorData) src.ElevatorData {	
	var buttonLights src.ElevatorData
	for floor := 0; floor < src.N_FLOORS; floor ++ {
		for direction := 0; direction < 2; direction ++ {
			for _, elevatorQueue := range elevatorQueues {
				if (elevatorQueue.OutsideOrders[floor][direction] == src.ORDER) {
					buttonLights.OutsideOrders[floor][direction] = src.ORDER
				}
			}
		}
	}

	buttonLights.InsideOrders = elevatorQueues[ourID].InsideOrders
	return buttonLights
}

func assignOrder(elevatorQueues map[string] src.ElevatorData, order src.ButtonOrder) {
	var minID string
	minCost := 100000

	for elevatorID, elevatorQueue := range elevatorQueues {
		if (elevatorQueue.OutsideOrders[order.Floor][order.ButtonType] == src.ORDER) {
			minID = elevatorID
			break;
		}

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
	cost := 0

	for floor := 0; floor < src.N_FLOORS; floor ++ {
		for direction := 0; direction < 2; direction ++ {
			if (elevator.OutsideOrders[floor][direction] == ORDER) {
				cost += COST_FOR_ORDER
			}

			if (elevator.InsideOrders[floor] == ORDER) {
				cost += COST_FOR_ORDER
			}
		}
	}

	cost+= COST_FOR_MOVEMENT * abs(elevator.Floor - order.Floor)
	return cost 
}


func calcNextOrderAndFloor(queueMatrix src.ElevatorData) src.ButtonOrder {
	currentOrder := src.ButtonOrder{-1, src.BUTTON_NONE}
	tmp := elevatorQueues[ourID]

	switch elevatorQueues[ourID].Direction {
		case src.DIR_UP, src.DIR_STOP:
			
			for floor := queueMatrix.Floor; floor < src.N_FLOORS - 1; floor ++ {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_UP] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_UP}
					tmp.Direction = src.DIR_UP
					elevatorQueues[ourID] = tmp
					return currentOrder
				}
			}

			for floor := src.N_FLOORS - 1; floor > 0; floor -- {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_DOWN] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_DOWN}
					if (floor < elevatorQueues[ourID].Floor) {
						tmp.Direction = src.DIR_DOWN
					} else {
						tmp.Direction = src.DIR_UP
					}
					elevatorQueues[ourID] = tmp
					return currentOrder
				}
			}

			for floor := 0; floor < queueMatrix.Floor; floor ++ {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_UP] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_UP}
					tmp.Direction = src.DIR_DOWN
					elevatorQueues[ourID] = tmp
					return currentOrder
				}
			}

		case src.DIR_DOWN:
			
			for floor := queueMatrix.Floor; floor > 0; floor -- {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_DOWN] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_DOWN}
					tmp.Direction = src.DIR_DOWN
					elevatorQueues[ourID] = tmp
					return currentOrder
				}
			}

			for floor := 0; floor < src.N_FLOORS - 1; floor ++ {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_UP] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_UP}
					if (floor > elevatorQueues[ourID].Floor) {
						tmp.Direction = src.DIR_UP
					} else {
						tmp.Direction = src.DIR_DOWN
					}
					elevatorQueues[ourID] = tmp
					return currentOrder
				}
			}

			for floor := src.N_FLOORS - 1; floor > queueMatrix.Floor; floor -- {
				if (queueMatrix.OutsideOrders[floor][src.BUTTON_DOWN] == ORDER || queueMatrix.InsideOrders[floor] == ORDER) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_DOWN}
					tmp.Direction = src.DIR_UP
					elevatorQueues[ourID] = tmp
					return currentOrder
				}
			}
		default:
			return currentOrder
	}

	tmp.Direction = src.DIR_STOP
	elevatorQueues[ourID] = tmp
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
	if (order.Floor != -1 && order.ButtonType != src.BUTTON_NONE) {
		tmp := elevatorQueues[ID]
		tmp.InsideOrders[order.Floor] = NO_ORDER
		tmp.OutsideOrders[order.Floor][order.ButtonType] = NO_ORDER
		elevatorQueues[ID] = tmp
	}
}

func timer(timeout chan bool) {
	for {
		time.Sleep(timeoutLimit)
		timeout <- true
	}
}


func abs(value int) int {
	if (value >= 0) { 	return value
	}else { 				return -1*value}
}


func commandPrint() {
	println(" ---- ORDERS ------")
	for id, queues := range elevatorQueues {
		println(id)
		tools.PrintQueue(queues)
	}
}


func InitQueue(chFloorFromController chan int, chOrderFromController chan src.ButtonOrder, chFinishedFromController chan bool, chGlobalOrdersToController chan src.ElevatorData, chDestinationFloorToController chan int) {
	var chUpdateGlobalOrders = make(chan bool)
 	go network.NetworkHandler()
	ourID = <- network.ChIDFromNetwork
 	var InitialQueue src.ElevatorData

 	go queueManager(chFloorFromController, chOrderFromController, chFinishedFromController,	chGlobalOrdersToController, chDestinationFloorToController,
					chUpdateGlobalOrders,  InitialQueue)
}


func queueManager(	chFloorFromController chan int,
					chOrderFromController chan src.ButtonOrder,
					chFinishedFromController chan bool,
					chGlobalOrdersToController chan src.ElevatorData,
					chDestinationFloorToController chan int,
					chUpdateGlobalOrders chan bool,
					InitialQueue src.ElevatorData) {	

	InitialQueue.Floor = <- chFloorFromController
	elevatorQueues[ourID] = InitialQueue
	network.ChQueueReadyToBeSent <- elevatorQueues[ourID]	
	go timer(chUpdateGlobalOrders)

	for {
		select {
			
			case <- chUpdateGlobalOrders:
				buttonLights = determineButtonLights(elevatorQueues)
				chGlobalOrdersToController <- buttonLights
				commandPrint()

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
				currentOrder = calcNextOrderAndFloor(elevatorQueues[ourID])
				if currentOrder.Floor != -1 {chDestinationFloorToController <- currentOrder.Floor}
				network.ChQueueReadyToBeSent <- elevatorQueues[ourID]

			case order := <- chOrderFromController:
				if (order.ButtonType != src.BUTTON_INSIDE){
					network.ChOrderFromQueue <- order
					assignOrder(elevatorQueues, order)
				} else {
					addOrder(ourID, order)
				}
				currentOrder = calcNextOrderAndFloor(elevatorQueues[ourID])
				if currentOrder.Floor != -1 {chDestinationFloorToController <- currentOrder.Floor}
				network.ChQueueReadyToBeSent <- elevatorQueues[ourID]

			case <- chFinishedFromController:
				deleteOrder(ourID, currentOrder)
				currentOrder = calcNextOrderAndFloor(elevatorQueues[ourID])
				if currentOrder.Floor != -1 {chDestinationFloorToController <- currentOrder.Floor}
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
		}
	}
}