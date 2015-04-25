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
var currentOrder	src.ButtonOrder
var allElevatorsData 	= make(map[string] src.ElevatorData)


func determineButtonLights(allElevatorsData map[string] src.ElevatorData) [src.N_FLOORS][3]int {	
	var buttonLights [src.N_FLOORS][3]int
	for floor := 0; floor < src.N_FLOORS; floor ++ {
		for buttonType := 0; buttonType < 2; buttonType ++ {
			for _, elevatorData := range allElevatorsData {
				if (elevatorData.OutsideOrders[floor][buttonType] == src.ORDER) {
					buttonLights[floor][buttonType] = src.ORDER
				}
			}
		}

		buttonLights[floor][2] = allElevatorsData[ourID].InsideOrders[floor]		
	}
	return buttonLights
}


func assignOrder(allElevatorsData map[string] src.ElevatorData, order src.ButtonOrder) {
	var minCostID string
	largeNumber := 100000
	minCost := largeNumber

	for elevatorID, elevatorData := range allElevatorsData {
		if (elevatorData.OutsideOrders[order.Floor][order.ButtonType] == src.ORDER) {
			minCostID = elevatorID
			break;
		}

		cost := calcOrderCost(elevatorData, order)
		if (cost < minCost) {
			minCostID = elevatorID
			minCost = cost
		
		} else if (cost == minCost && elevatorID < minCostID) {
			minCostID = elevatorID
			minCost = cost
		}
	}

	if (minCostID == ourID) {
		addOrder(ourID, order)
		
	}
}


func calcOrderCost(elevatorData src.ElevatorData, order src.ButtonOrder) int {
	cost := 0

	for floor := 0; floor < src.N_FLOORS; floor ++ {
		for buttonType := 0; buttonType < 2; buttonType ++ {
			if (elevatorData.OutsideOrders[floor][buttonType] == ORDER) {
				cost += COST_FOR_ORDER
			}

			if (elevatorData.InsideOrders[floor] == ORDER) {
				cost += COST_FOR_ORDER
			}
		}
	}

	cost+= COST_FOR_MOVEMENT * abs(elevatorData.Floor - order.Floor)
	return cost 
}


func isOrderAtFloor(elevatorData src.ElevatorData, buttonType int, floor int) bool {
	return elevatorData.OutsideOrders[floor][buttonType] == ORDER || elevatorData.InsideOrders[floor] == ORDER
}

func calcDirection(currentFloor int, destinationFloor int, direction int) int{
	switch direction {
		
		case src.DIR_DOWN:
			
			if (destinationFloor > currentFloor) {
				return src.DIR_UP
			} else {
				return src.DIR_DOWN
			}

		case src.DIR_UP, src.DIR_STOP:
			
			if (destinationFloor < currentFloor) {
				return src.DIR_DOWN
			} else {
				return src.DIR_UP
			}
	}
	return direction
}

func calcNextOrderAndFloor(elevatorData src.ElevatorData) src.ButtonOrder {
	currentOrder := src.ButtonOrder{-1, src.BUTTON_NONE}
	tmp := allElevatorsData[ourID]

	switch allElevatorsData[ourID].Direction {
		case src.DIR_UP, src.DIR_STOP:
			
			for floor := elevatorData.Floor; floor < src.N_FLOORS - 1; floor ++ {
				if (isOrderAtFloor(elevatorData, src.BUTTON_UP, floor)) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_UP}
					tmp.Direction = calcDirection(allElevatorsData[ourID].Floor, floor, tmp.Direction)
					allElevatorsData[ourID] = tmp
					return currentOrder
				}
			}

			for floor := src.N_FLOORS - 1; floor > 0; floor -- {
				if (isOrderAtFloor(elevatorData, src.BUTTON_DOWN, floor)) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_DOWN}
					tmp.Direction = calcDirection(allElevatorsData[ourID].Floor, floor, tmp.Direction)
					allElevatorsData[ourID] = tmp
					return currentOrder
				}
			}

			for floor := 0; floor < elevatorData.Floor; floor ++ {
				if (isOrderAtFloor(elevatorData, src.BUTTON_UP, floor)) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_UP}
					tmp.Direction = calcDirection(allElevatorsData[ourID].Floor, floor, tmp.Direction)
					allElevatorsData[ourID] = tmp
					return currentOrder
				}
			}

		case src.DIR_DOWN:
			
			for floor := elevatorData.Floor; floor > 0; floor -- {
				if (isOrderAtFloor(elevatorData, src.BUTTON_DOWN, floor)) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_DOWN}
					tmp.Direction = calcDirection(allElevatorsData[ourID].Floor, floor, tmp.Direction)
					allElevatorsData[ourID] = tmp
					return currentOrder
				}
			}

			for floor := 0; floor < src.N_FLOORS - 1; floor ++ {
				if (isOrderAtFloor(elevatorData, src.BUTTON_UP, floor)) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_UP}
					tmp.Direction = calcDirection(allElevatorsData[ourID].Floor, floor, tmp.Direction)
					allElevatorsData[ourID] = tmp
					return currentOrder
				}
			}

			for floor := src.N_FLOORS - 1; floor > elevatorData.Floor; floor -- {
				if (isOrderAtFloor(elevatorData, src.BUTTON_DOWN, floor)) {
					currentOrder = src.ButtonOrder{floor, src.BUTTON_DOWN}
					tmp.Direction = calcDirection(allElevatorsData[ourID].Floor, floor, tmp.Direction)
					allElevatorsData[ourID] = tmp
					return currentOrder
				}
			}
		default:
			return currentOrder
	}

	tmp.Direction = src.DIR_STOP
	allElevatorsData[ourID] = tmp
	return currentOrder
}


func addOrder(ID string, order src.ButtonOrder) {
	switch(order.ButtonType) {

		case src.BUTTON_INSIDE:
			tmp := allElevatorsData[ID]
			tmp.InsideOrders[order.Floor] = ORDER
			allElevatorsData[ID] = tmp
		
		default:
			tmp := allElevatorsData[ID]
			tmp.OutsideOrders[order.Floor][order.ButtonType] = ORDER
			allElevatorsData[ID] = tmp
	}
}


func deleteOrder(ID string, order src.ButtonOrder) {
	if (order.Floor != -1 && order.ButtonType != src.BUTTON_NONE) {
		tmp := allElevatorsData[ID]
		tmp.InsideOrders[order.Floor] = NO_ORDER
		tmp.OutsideOrders[order.Floor][order.ButtonType] = NO_ORDER
		allElevatorsData[ID] = tmp
	}
}

func updateLightsTimer(timeout chan bool) {
	updateLightsInterval := 100*time.Millisecond
	
	for {
		time.Sleep(updateLightsInterval)
		timeout <- true
	}
}


func abs(value int) int {
	if (value >= 0) { 	
		return value
	}else { 				
		return -1*value
	}
}


func commandPrint() {
	println(" ---- ORDERS ------")
	for id, queues := range allElevatorsData {
		println(id)
		tools.PrintQueue(queues)
	}
}


func InitQueue(chFloorFC chan int, chOrderFC chan src.ButtonOrder, chOrderFinishedFC chan bool, chButtonLightsTC chan [src.N_FLOORS][3]int, chDestinationFloorTC chan int) {
	var chUpdateLights = make(chan bool)
	var chIdentificationFN = make(chan string)
	var chElevatorDataFN = make(chan network.ElevatorMessage)
	var chElevatorDataTN = make(chan src.ElevatorData, 2)
	var chOrderFN = make(chan src.ButtonOrder)
	var chOrderTN = make(chan src.ButtonOrder, 2)
	var chDisconElevatorFN = make(chan string)

 	go network.InitNetwork(chIdentificationFN, chElevatorDataFN, chElevatorDataTN, chOrderFN, chOrderTN, chDisconElevatorFN)

	ourID = <- chIdentificationFN
 	var InitialElevatorData src.ElevatorData

 	InitialElevatorData.Floor = <- chFloorFC
	allElevatorsData[ourID] = InitialElevatorData

	chElevatorDataTN <- allElevatorsData[ourID]	
	go updateLightsTimer(chUpdateLights)

 	go queueManager(chFloorFC, chOrderFC, chOrderFinishedFC, chButtonLightsTC, chDestinationFloorTC, chUpdateLights,
 					chIdentificationFN, chElevatorDataFN, chElevatorDataTN, chOrderFN, chOrderTN, chDisconElevatorFN)
}


func queueManager(	chFloorFC chan int,
					chOrderFC chan src.ButtonOrder,
					chOrderFinishedFC chan bool,
					chButtonLightsTC chan [src.N_FLOORS][3]int,
					chDestinationFloorTC chan int,
					chUpdateLights chan bool,
					chIdentificationFN chan string,
					chElevatorDataFN chan network.ElevatorMessage,
					chElevatorDataTN chan src.ElevatorData,
					chOrderFN chan src.ButtonOrder,
					chOrderTN chan src.ButtonOrder,
					chDisconElevatorFN chan string) {
	for {
		select {
			case <- chUpdateLights:
				buttonLights := determineButtonLights(allElevatorsData)
				chButtonLightsTC <- buttonLights
				commandPrint()

			case order := <- chOrderFN:
				assignOrder(allElevatorsData, order)
				currentOrder = calcNextOrderAndFloor(allElevatorsData[ourID])
				if currentOrder.Floor != -1 {chDestinationFloorTC <- currentOrder.Floor}
				chElevatorDataTN <- allElevatorsData[ourID]

			case elevatorMessage := <- chElevatorDataFN:
				allElevatorsData[elevatorMessage.SenderAddress] = elevatorMessage.ElevatorData

			case floor := <- chFloorFC:
				tmp := allElevatorsData[ourID]
				tmp.Floor = floor
				allElevatorsData[ourID] = tmp
				currentOrder = calcNextOrderAndFloor(allElevatorsData[ourID])
				if currentOrder.Floor != -1 {chDestinationFloorTC <- currentOrder.Floor}
				chElevatorDataTN <- allElevatorsData[ourID]

			case order := <- chOrderFC:
				if (order.ButtonType != src.BUTTON_INSIDE){
					chOrderTN <- order
					assignOrder(allElevatorsData, order)
				} else {
					addOrder(ourID, order)
				}
				currentOrder = calcNextOrderAndFloor(allElevatorsData[ourID])
				if currentOrder.Floor != -1 {chDestinationFloorTC <- currentOrder.Floor}
				chElevatorDataTN <- allElevatorsData[ourID]

			case <- chOrderFinishedFC:
				deleteOrder(ourID, currentOrder)
				currentOrder = calcNextOrderAndFloor(allElevatorsData[ourID])
				if currentOrder.Floor != -1 {chDestinationFloorTC <- currentOrder.Floor}
				chElevatorDataTN <- allElevatorsData[ourID]

			case elevator := <- chDisconElevatorFN:
				ordersToDistribute := allElevatorsData[elevator]
				delete(allElevatorsData, elevator)
				
				for floor := 0; floor < src.N_FLOORS; floor ++ {
					for buttonType := 0; buttonType < 2; buttonType ++ {
						if (ordersToDistribute.OutsideOrders[floor][buttonType] == src.ORDER) {
							assignOrder(allElevatorsData, src.ButtonOrder{floor, buttonType})
						}
					}
				}

				chElevatorDataTN <- allElevatorsData[ourID]
				currentOrder = calcNextOrderAndFloor(allElevatorsData[ourID])
				if currentOrder.Floor != -1 {chDestinationFloorTC <- currentOrder.Floor}
		}
	}
}