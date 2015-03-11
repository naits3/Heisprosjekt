

package Queue//package queue

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

//Definerer egne variabler:
var localQueue     	elevatorData
var globalQueues  []elevatorData

const FLOORS = 4 //Import from Styring!


func MergeOrders(queueList []elevatorData) elevatorData {
	//Go through list, append orders. Remove if orders are finished
	//Remember to add your own global list as well
	var mergedData elevatorData
	var mergedQueue [FLOORS][2]int	

	for floor := 0; floor < FLOORS; floor ++ {		
		for direction := 0; direction < 2; direction ++ {
			
			for eachQueue := 0; eachQueue < len(queueList); eachQueue ++ {
				switch queueList[eachQueue].outsideOrders[floor][direction] {
					case ORDER:
						mergedQueue[floor][direction] = ORDER
					case DELETE_ORDER:
						mergedQueue[floor][direction] = EMPTY
						break;
				}
			}
		}
	}

	mergedData.outsideOrders = mergedQueue
	return mergedData


}

func assignOrders(queueList []elevatorData, mergedQueue elevatorData) int { // return elevatorData
	//queueList = globalQueues + localQueue
	//for direction = 0; direction < 2; direction ++:
		//for floor = 0; floor < FLOOR; floor ++:
			//if mergedQueue.outsideOrders[floor][direction] == 1:
				//Altsaa, hvis bestilling ...
				//for elevator = 0; elevator < len(queueList); elevator ++:
					//calcTotalCost(elevator)
					//minst totalCost faar ordren

	return 0			

}

func calcTotalCost() int {
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
