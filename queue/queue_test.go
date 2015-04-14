package queue

import (
	"testing"
	"time"
	"Heisprosjekt/src"
	"Heisprosjekt/tools"
	//"Heisprosjekt/network"
)


func TestMergeOrders(t *testing.T) {
	var testData1 src.ElevatorData
	// var testData2 src.ElevatorData
	// var testData3 src.ElevatorData
	//var mergedData elevatorData
	
	testData1.OutsideOrders[0][1] = ORDER
	testData1.OutsideOrders[2][0] = DELETE_ORDER
	testData1.OutsideOrders[2][1] = ORDER
	testData1.InsideOrders[2] = DELETE_ORDER

	queueList := []src.ElevatorData{testData1}
	
	// ------- ONLY FOR PRINTING: -------------//
	println("Elevator 1:")
	tools.PrintQueue(testData1)
	// println("Elevator 2:")
	// tools.PrintQueue(testData2)
	// println("Elevator 3:")
	// tools.PrintQueue(testData3)
	println(" ---------------- ")
	// ---------------------------------------//

	mergedData := mergeOrders(queueList)
	println("After Merge:")
	tools.PrintQueue(mergedData)
}


func TestCalcTotalCost(t *testing.T) {
	var testData1 src.ElevatorData
	testData1.Direction = src.DIR_STOP
	testData1.Floor = 1
	//testData1.OutsideOrders[1][1] = ORDER
	//testData1.OutsideOrders[1][0] = ORDER
	testData1.OutsideOrders[0][0] = ORDER
	testData1.OutsideOrders[1][1] = ORDER
	//testData1.insideOrders[1] = ORDER

	tools.PrintQueue(testData1)

	cost := calcTotalCost(&testData1)
	println("Cost: ",cost)

	tools.PrintQueue(testData1)
	// expect := 13
	// if cost != expect {
	// 	t.Error("expected", expect, "got", cost)
	// } 
}


func TestFindMinimum(t *testing.T){
	//array := []int{8,5,3,4}
	// println("Lowest index is: ",findMinimumCost(array))
}


func TestClearOutsideOrders(t *testing.T) {
	var testData1 src.ElevatorData
	var testData2 src.ElevatorData
	var testData3 src.ElevatorData
	testData1.OutsideOrders[0][0] = ORDER
	testData1.InsideOrders[2] = ORDER
	testData1.ID = 1
	testData2.OutsideOrders[1][0] = ORDER
	testData2.InsideOrders[2] = ORDER
	testData2.InsideOrders[3] = ORDER
	testData2.ID = 2
	testData3.OutsideOrders[2][1] = ORDER
	testData3.ID = 3

	dataArray := []src.ElevatorData{testData1, testData2, testData3}
	//dataArray = []src.ElevatorData{}

	println("Before clearing: ")
	for elevator := 0; elevator < len(dataArray); elevator ++ {
		tools.PrintQueue(dataArray[elevator])
	}

	println("After clearing: ")
	clearOutsideOrders(dataArray)

	for elevator := 0; elevator < len(dataArray); elevator ++ {
		tools.PrintQueue(dataArray[elevator])
	}
}


func TestAssignOrders(t *testing.T) {
	var mergedQueue src.ElevatorData
	var elevatorOne src.ElevatorData
	var elevatorTwo src.ElevatorData
	var elevatorThree src.ElevatorData

	mergedQueue.OutsideOrders[2][1] = ORDER
	mergedQueue.OutsideOrders[1][0] = ORDER
	mergedQueue.OutsideOrders[3][0] = ORDER
	mergedQueue.OutsideOrders[2][0] = ORDER

	elevatorOne.ID = 7
	elevatorOne.Direction = src.DIR_STOP
	elevatorOne.Floor = 2
	elevatorOne.InsideOrders[1] = ORDER	
	elevatorOne.OutsideOrders[1][0] = ORDER
	elevatorOne.OutsideOrders[2][1] = ORDER

	elevatorTwo.ID = 8
	elevatorTwo.Direction = src.DIR_STOP
	elevatorTwo.Floor = 0
	elevatorTwo.InsideOrders[2] = ORDER
	elevatorTwo.OutsideOrders[3][1] = ORDER

	elevatorThree.ID = 9
	elevatorThree.Direction = src.DIR_STOP
	elevatorThree.Floor = 2
	elevatorThree.InsideOrders[1] = ORDER

	println("Before assignment: ")
	tools.PrintQueue(mergedQueue)

	// queueArray2 := []src.ElevatorData{elevatorThree, elevatorOne, elevatorTwo}
	// queueArray1 := []src.ElevatorData{elevatorTwo, elevatorThree, elevatorOne}
	// queueArray3 := []src.ElevatorData{elevatorTwo, elevatorOne, elevatorThree}


	// assignedOrder1 := assignOrders(queueArray1, mergedQueue)
	// assignedOrder2 := assignOrders(queueArray2, mergedQueue)
	// assignedOrder3 := assignOrders(queueArray3, mergedQueue)

	// println("After assignment: ")
}


func TestQueueHandler(t *testing.T) {

	var done					= make(chan bool)
	var chNewFloor				= make(chan int)
	var chNewOrder 				= make(chan src.ButtonOrder)
	var chNewDirection			= make(chan int)
	var chOrderIsFinished 		= make(chan bool)
	var chNewOrdersFromQueue 	= make(chan src.ElevatorData)
	var chNewNextFloorFromQueue = make(chan int)
	
	knownOrders.OutsideOrders[0][0] = 1
	knownOrders.OutsideOrders[2][1] = 1
	knownOrders.OutsideOrders[2][0] = 1
	knownOrders.InsideOrders[0] = 1


	InitQueue(chNewFloor, chNewOrder, chNewDirection, chOrderIsFinished, chNewOrdersFromQueue, chNewNextFloorFromQueue)


	go func() {
		for {
			select {
				case <- chNewOrdersFromQueue:
					continue
				case <- chNewNextFloorFromQueue:
					continue
				default:
					time.Sleep(100*time.Millisecond)

			}
		}
	}()

	time.Sleep(10*time.Second)

	// chNewOrder <- src.ButtonOrder{1, src.BUTTON_UP}
	// knownOrders.OutsideOrders[2][0] = DELETE_ORDER

	// chNewDirection <- src.DIR_STOP

	// time.Sleep(10*time.Second)

	chOrderIsFinished <- true

	<- done
}

func TestDebug(t *testing.T) {
	//knownOrders.OutsideOrders[2][1] = 1
	//knownOrders.OutsideOrders[0][0] = 1
	//knownOrders.OutsideOrders[3][0] = 1
	//knownOrders.InsideOrders[3] = 1
	knownOrders.InsideOrders[1] = 1
	knownOrders.ID = 3
	knownOrders.Floor = 3

	allElevatorData := []src.ElevatorData{knownOrders}
	mergedQueue := mergeOrders(allElevatorData)
	knownOrders.OutsideOrders = mergedQueue.OutsideOrders

	assignedOrder := assignOrders(allElevatorData, mergedQueue)
	currentOrder = calcNextOrderAndFloor(assignedOrder)

	println("Assigned Order: ")
	tools.PrintQueue(assignedOrder)
	println("Next floor: ", currentOrder.Floor)
	println("Button Type: ", currentOrder.ButtonType)
}


