

package Queue//package Queue

type elevatorData struct { 
	floor 			int
	direction 		int
	outsideOrders [FLOORS][2]int
	insideOrders  []int
}

type Order struct {	// THis struct should be imported from styring / IO!
	floor 			int
	buttonType		int
}	

const (
	ORDER = 1
	EMPTY = 0
	DELETE_ORDER = -1
)

const (
	UP = 1
	IDLE = 0
	DOWN = -1
)


const FLOORS = 4 //Import from Styring!
var localQueue     	elevatorData
var globalQueues  []elevatorData


func MergeOrders(queueList []elevatorData) elevatorData {
	//Go through list, append orders. Remove if orders are finished
	//Remember to add your own global list as well
	var mergedData elevatorData
	var mergedQueue [FLOORS][2]int	

	for floor := 0; floor < FLOORS; floor ++ {
		directionLoop:	
		for direction := 0; direction < 2; direction ++ {
			
			for eachQueue := 0; eachQueue < len(queueList); eachQueue ++ {
				switch queueList[eachQueue].outsideOrders[floor][direction] {
					case ORDER:
						mergedQueue[floor][direction] = ORDER
					case DELETE_ORDER:
						mergedQueue[floor][direction] = EMPTY
						break directionLoop
				}
			}
		}
	}

	mergedData.outsideOrders = mergedQueue
	return mergedData
}

func assignOrders(queueList []elevatorData, mergedQueue elevatorData) int { // return elevatorData
	//queueList = globalQueues + localQueue
	for floor := 0; floor < FLOORS; floor ++ {
		for direction := 0; direction < 2; direction ++ {
			
			if mergedQueue.outsideOrders[floor][direction] == ORDER {
				for eachElevator := 0; eachElevator < len(queueList); eachElevator ++{
					calcTotalCost(queueList[eachElevator])
				}
				//minst totalCost faar ordren
				//Hvilke tall skal vi bruke for aa gi unik ID til hver heis?
				// 1 er jo ORDER!
			}
		}
	}

	return 0			

}

func calcTotalCost(queueData elevatorData) int {
	const COST_FOR_ORDER = 3
	const COST_FOR_MOVE = 1

	switch queueData.direction {
		case UP:
			// ...
		case DOWN:
			// ...
		default: //Ta med IDLE?
			// ...
	}

	// Gaa gjennom hele koen (bruk dir for aa bestemme start-retning) og legg til costs
	// for hver etasje/ordre helt til alle ordrene er tatt med!
	return 0;
}

func calcNextFloor() {

}


func main() { // func queueHandler() {
	// Initialisere kanaler her:
	// chReadyToMerge 	:= make(chan bool)
	// chIncomingQueues:= make(chan elevatorData)
	// chNewFloor		:= make(chan int)
	// chNewOrder 		:= make(chan Order) //make(chan styring.Order)
	// chNewDirection	:= make(chan int)
	// chDeleteOrder 	:= make(chan Order) //make(chan styring.Order)

	// for {
	// 	select {
	// 		case <- chReadyToMerge:
	// 			//add our elevatorData into dataStorage
	// 			//mergeOrders()
	// 			//assignOrders()
	// 			//calcNextFloor()
	// 			//send mergedOrders til nettverk
	// 			//toem dataStorage
	// 		case data := <- chIncomingQueues:
	// 			//append to list of data
	// 		case newFloor := <- chNewFloor:
	// 			//update elevatorData.floor
	// 		case newOrder := <- chNewOrder:
	// 			//update elevatorData.queueMatrix
	// 		case deletedOrder := <- chDeleteOrder:
	// 			//update elevatorData.queueMatrix with a special char
	// 		case newDirection := <- chNewDirection:
	// 			//update elevator.direction

	// 	}
	// }
}
