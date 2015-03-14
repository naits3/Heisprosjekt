

package queue//package Queue

import (
	"Heisprosjekt/src"
	"Heisprosjekt/network"
)

const (
	ORDER = 1
	EMPTY = 0
	DELETE_ORDER = -1
)

var localQueue     	src.ElevatorData
var listOfIncomingData  []src.ElevatorData

// Channels made for communication across modules



func mergeOrders(queueList []src.ElevatorData) src.ElevatorData {
	//Go through list, append orders. Remove if orders are finished
	//Remember to add your own global list as well
	var mergedData src.ElevatorData
	var mergedQueue [src.N_FLOORS][2]int	

	for floor := 0; floor < src.N_FLOORS; floor ++ {
		directionLoop:	
		for direction := 0; direction < 2; direction ++ {
			
			for eachQueue := 0; eachQueue < len(queueList); eachQueue ++ {
				switch queueList[eachQueue].OutsideOrders[floor][direction] {
					case ORDER:
						mergedQueue[floor][direction] = ORDER
					case DELETE_ORDER:
						mergedQueue[floor][direction] = EMPTY
						break directionLoop
				}
			}
		}
	}

	mergedData.OutsideOrders = mergedQueue
	return mergedData
}

func assignOrders(queueList []src.ElevatorData, mergedQueue src.ElevatorData) int { // return elevatorData
	var costArray []int
	//queueList = globalQueues + localQueue
	for floor := 0; floor < src.N_FLOORS; floor ++ {
		for direction := 0; direction < 2; direction ++ {
			
			if mergedQueue.OutsideOrders[floor][direction] == ORDER {
				for eachElevator := 0; eachElevator < len(queueList); eachElevator ++{
					costArray = append(costArray, calcTotalCost(queueList[eachElevator]))
				}
				//minst totalCost faar ordren
				//Hvilke tall skal vi bruke for aa gi unik ID til hver heis? Trenger vel ikke det?
				// 1 er jo ORDER!
			}
		}
	}

	return 0			

}

func calcTotalCost(queueData src.ElevatorData) int {
	const COST_FOR_ORDER = 3
	const COST_FOR_MOVE = 1
	totalCost := 0
	floorsSinceLastOrder := 0

	// remove "illegal states" here?
	// i.e dir = DIR_UP & floor = 4
	// and dir = src.DIR_DOWN & floor = 0

	// IDE:
	// Kan man gjoere det saa enkelt at man teller oppover og nedover helt til man finner siste 1-er
	// Deretter teller antall 1-ere?
	switch queueData.Direction {
		case src.DIR_UP:

			for floor := queueData.Floor; floor < src.N_FLOORS; floor ++ {
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1
				if queueData.OutsideOrders[floor][0] == ORDER || queueData.InsideOrders[floor] == ORDER {
					totalCost += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}				
			}

			floorsSinceLastOrder -= 1 

			for floor := src.N_FLOORS-1; floor >= 0; floor-- {
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1 
				if queueData.OutsideOrders[floor][1] == ORDER || (queueData.InsideOrders[floor] == ORDER  && floor < queueData.Floor){
					totalCost += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}				
			}
			
			floorsSinceLastOrder -= 1

			for floor := 0; floor < queueData.Floor; floor ++{
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1
				if queueData.OutsideOrders[floor][0] == ORDER {
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
				if queueData.OutsideOrders[floor][1] == ORDER || queueData.InsideOrders[floor] == ORDER  {
					totalCost += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}				
			}
 
			floorsSinceLastOrder -= 1 

			for floor := 0; floor < src.N_FLOORS; floor ++ {
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1
				if queueData.OutsideOrders[floor][0] == ORDER || (queueData.InsideOrders[floor] == ORDER && floor > queueData.Floor){
					totalCost += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}				
			}

			floorsSinceLastOrder -= 1

			for floor := src.N_FLOORS-1; floor > queueData.Floor; floor --{
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1
				if queueData.OutsideOrders[floor][1] == ORDER {
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


func QueueHandler() { // func queueHandler() {
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
