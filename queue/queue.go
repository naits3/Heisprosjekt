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

const (
	COST_FOR_ORDER = 2
	COST_FOR_MOVEMENT = 1

)

var ourID			string
var globalOrders    src.ElevatorData
var currentOrder	src.ButtonOrder
var elevatorQueues	= make(map[string] src.ElevatorData)
var timeoutLimit time.Duration = 500*time.Millisecond
var queueDirection	int


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

// TO BE EDITED!
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

func QueueManager(chFloorFromController chan int, chOrderFromController chan src.ButtonOrder, chFinishedFromController chan bool, chGlobalOrdersToController chan src.ElevatorData, chDestinationFloorToController chan int) {
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
				for id, queues := range elevatorQueues {
					println(id)
					tools.PrintQueue(queues)
				}

			case order := <- network.ChOrderToQueue:
				println("Got order from NETwork!")
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

			//case direction := <- chDirectionFromController:
			//	println("Q: direction = ",direction)
			//	tmp := elevatorQueues[ourID]
			//	tmp.Direction = direction
			//	elevatorQueues[ourID] = tmp
			//	network.ChQueueReadyToBeSent <- elevatorQueues[ourID]

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
// * renew calcNextFloor and add functionality for DIR_STOP(behøver vi egentlig det?)

// BUGS:
// * Hvis vi står i 2. etg, og behandler "ned" ordre, og dermed trykker på "ned"-ordre i 3. etg, vil 
//   heisen slette ordren i 3... Dette er fordi calcNextFloor tror vi står i idle, og beregner basert 
//   på at vi er på vei opp																				| FIXED

// * Heisen tar alle opp- og ned- etasjer når den går ovenfra og nedover. Dette er fordi calc cost 
//   beregner likt for opp og idle, og når en heis står i en etasje er den i idle. 						| FIXED

// * Dersom heiser står i 3., og får en bestilling i 1. går den ned. Når den er mellom 3. og 2., og den 
//   får en ny bestilling i 3, vil den TRO at det står i 3. og sier at den skal ta denne. Men den skal 
//   jo ta bestillingen i 2. (Ikke 1.pri å fikse)

// * Heisen mottar bestillinger over nettet, men det er ikke alltid at den tar bestillingen.