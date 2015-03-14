

package queue//package Queue

type elevatorData struct { 
	IP 				int
	floor 			int
	direction 		int
	outsideOrders [FLOORS][2]int
	insideOrders  [FLOORS]int
}

type Order struct {	// THis struct should be imported from styring 
	floor 			int
	buttonType		int
}	

const (
	ORDER = 1
	EMPTY = 0
	DELETE_ORDER = -1
)

const ( //Importer fra styring
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

func AssignOrders(queueList []elevatorData, mergedQueue elevatorData) int { // return elevatorData
	var costArray = []int
	//queueList = globalQueues + localQueue
	for floor := 0; floor < FLOORS; floor ++ {
		for direction := 0; direction < 2; direction ++ {
			
			if mergedQueue.outsideOrders[floor][direction] == ORDER {
				for eachElevator := 0; eachElevator < len(queueList); eachElevator ++{
					costArray = append(costArray, CalcTotalCost(queueList[eachElevator]))
				}
				//minst totalCost faar ordren
				//Hvilke tall skal vi bruke for aa gi unik ID til hver heis? Trenger vel ikke det?
				// 1 er jo ORDER!
			}
		}
	}

	return 0			

}

func CalcTotalCost(queueData elevatorData) int {
	const COST_FOR_ORDER = 3
	const COST_FOR_MOVE = 1
	totalCost := 0
	floorsSinceLastOrder := 0

	// remove "illegal states" here?
	// i.e dir = UP & floor = 4
	// and dir = DOWN & floor = 0

	// IDE:
	// Kan man gjoere det saa enkelt at man teller oppover og nedover helt til man finner siste 1-er
	// Deretter teller antall 1-ere?
	switch queueData.direction {
		case UP:

			for floor := queueData.floor; floor < FLOORS; floor ++ {
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1
				if queueData.outsideOrders[floor][0] == (ORDER || queueData.IP) || queueData.insideOrders[floor] == (ORDER || queueData.IP){
					totalCost += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}				
			}

			floorsSinceLastOrder -= 1 

			for floor := FLOORS-1; floor >= 0; floor-- {
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1 
				if queueData.outsideOrders[floor][1] == (ORDER || queueData.IP) || (queueData.insideOrders[floor] == (ORDER || queueData.IP) && floor < queueData.floor){
					totalCost += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}				
			}
			
			floorsSinceLastOrder -= 1

			for floor := 0; floor < queueData.floor; floor ++{
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1
				if queueData.outsideOrders[floor][0] == (ORDER || queueData.IP) {
					totalCost  += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}
			}

			totalCost -= 3*COST_FOR_MOVE // 1 when we start, and 1 when we change dir * 2
			totalCost -= floorsSinceLastOrder
			return totalCost
		
		case DOWN:

			for floor := queueData.floor; floor >= 0; floor-- {
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1 
				if queueData.outsideOrders[floor][1] == (ORDER || queueData.IP) || queueData.insideOrders[floor] == (ORDER || queueData.IP) {
					totalCost += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}				
			}
 
			floorsSinceLastOrder -= 1 

			for floor := 0; floor < FLOORS; floor ++ {
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1
				if queueData.outsideOrders[floor][0] == (ORDER || queueData.IP) || (queueData.insideOrders[floor] == (ORDER || queueData.IP) && floor > queueData.floor){
					totalCost += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}				
			}

			floorsSinceLastOrder -= 1

			for floor := FLOORS-1; floor > queueData.floor; floor --{
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1
				if queueData.outsideOrders[floor][1] == (ORDER || queueData.IP) {
					totalCost  += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}
			}

			totalCost -= 3*COST_FOR_MOVE // 1 when we start, and 1 when we change dir * 2
			totalCost -= floorsSinceLastOrder
			return totalCost	

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
