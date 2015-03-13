

package Queue//package Queue

type elevatorData struct { 
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
	//queueList = globalQueues + localQueue
	for floor := 0; floor < FLOORS; floor ++ {
		for direction := 0; direction < 2; direction ++ {
			
			if mergedQueue.outsideOrders[floor][direction] == ORDER {
				for eachElevator := 0; eachElevator < len(queueList); eachElevator ++{
					CalcTotalCost(queueList[eachElevator])
				}
				//minst totalCost faar ordren
				//Hvilke tall skal vi bruke for aa gi unik ID til hver heis?
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
				if queueData.outsideOrders[floor][0] == ORDER || queueData.insideOrders[floor] == ORDER{
					totalCost += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}				
			}

			floorsSinceLastOrder -= 1 

			for floor := FLOORS-1; floor >= 0; floor-- {
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1 
				if queueData.outsideOrders[floor][1] == ORDER || (queueData.insideOrders[floor] == ORDER && floor < queueData.floor){
					totalCost += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}				
			}
			
			floorsSinceLastOrder -= 1

			for floor := 0; floor < queueData.floor; floor ++{
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1
				if queueData.outsideOrders[floor][0] == ORDER {
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
				if queueData.outsideOrders[floor][1] == ORDER || queueData.insideOrders[floor] == ORDER {
					totalCost += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}				
			}
 
			floorsSinceLastOrder -= 1 

			for floor := 0; floor < FLOORS; floor ++ {
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1
				if queueData.outsideOrders[floor][0] == ORDER || (queueData.insideOrders[floor] == ORDER && floor > queueData.floor){
					totalCost += COST_FOR_ORDER
					floorsSinceLastOrder = 0
				}				
			}

			floorsSinceLastOrder -= 1

			for floor := FLOORS-1; floor > queueData.floor; floor --{
				totalCost += COST_FOR_MOVE
				floorsSinceLastOrder += 1
				if queueData.outsideOrders[floor][1] == ORDER {
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

func CalcTotalCostv2(queueData elevatorData) int {
	// Proever ut idéen om aa telle oppover eller nedover helt til man finner siste 1-er, og deretter legger til
	// antall 1-ere
	const COST_FOR_ORDER = 3
	const COST_FOR_MOVE = 1
	totalCost := 0
	//floorsSinceLastOrder := 0

	switch queueData.direction {
		case UP:

			// PROCEDURE:
			// RUN THROUGH ALL THE FLOORS UPWARDS, TURN AT THE TOP FLOOR AND RUN DOWNWARDS UNTIL 1ST FLOOR,
			// RUN UPWARDS AGAIN UNTIL YOU REACH THE FLOOR BENEATH YOUR CURRENT FLOOR. COUNT #MOVEMENTS
			// TO THE LAST ORDER

			// for floor := queueData.floor; floor < FLOORS; floor ++ {
			// 	totalCost += COST_FOR_MOVE // saa var det det problemet her igjen, da.. 
			// }
			// ...
			break
		case DOWN:
			// ...
			break
		default:
			// ...
	}
	// Gaar det an aa legge forloop her, slik at case-ene kun setter hvordan man skal starte aa iterere?
	// RUN THROUGH THE QUEUE AGAIN, AND COUNT ALL THE ORDERS. APPEND THE TOTALCOST

	return totalCost;

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
