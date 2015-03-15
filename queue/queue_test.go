package queue

import (
	"testing"
	"Heisprosjekt/src"
	"Heisprosjekt/tools"
)


func TestMergeOrders(t *testing.T) {
	var testData1 src.ElevatorData
	var testData2 src.ElevatorData
	var testData3 src.ElevatorData
	//var mergedData elevatorData
	
	testData1.OutsideOrders[0][0] = ORDER
	testData2.OutsideOrders[2][1] = DELETE_ORDER
	testData3.OutsideOrders[2][1] = DELETE_ORDER
	testData3.InsideOrders[1] = ORDER

	queueList := []src.ElevatorData{testData1, testData2, testData3}
	
	// ------- ONLY FOR PRINTING: -------------//
	println("Elevator 1:")
	tools.PrintQueue(testData1)
	println("Elevator 2:")
	tools.PrintQueue(testData2)
	println("Elevator 3:")
	tools.PrintQueue(testData3)
	println(" ---------------- ")
	// ---------------------------------------//

	mergedData := mergeOrders(queueList)
	println("After Merge:")
	tools.PrintQueue(mergedData)
}


func TestCalcTotalCost(t *testing.T) {
	var testData1 src.ElevatorData
	testData1.Direction = src.DIR_UP
	testData1.Floor = 1
	testData1.OutsideOrders[2][1] = ORDER
	//testData1.OutsideOrders[1][1] = ORDER
	testData1.OutsideOrders[1][0] = ORDER
	//testData1.insideOrders[1] = ORDER

	tools.PrintQueue(testData1)

	cost := calcTotalCost(testData1)
	println(cost)

	expect := 13
	if cost != expect {
		t.Error("expected", expect, "got", cost)
	} 
}

func TestFindMinimum(t *testing.T){
	//array := []int{8,5,3,4}
	// println("Lowest index is: ",findMinimumCost(array))
}

func TestClearOutsideOrders(t *testing.T) {
	var testData1 src.ElevatorData
	var testData2 src.ElevatorData
	var testData3 src.ElevatorData
	var ourQueue  src.ElevatorData
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
	clearOutsideOrders(dataArray, ourQueue)

	for elevator := 0; elevator < len(dataArray); elevator ++ {
		tools.PrintQueue(dataArray[elevator])
	}
}

func TestAssignOrders(t *testing.T) {
	var mergedQueue src.ElevatorData
	var elevatorOne src.ElevatorData
	var elevatorTwo src.ElevatorData

	mergedQueue.OutsideOrders[2][1] = ORDER
	mergedQueue.OutsideOrders[1][0] = ORDER
	mergedQueue.OutsideOrders[3][0] = ORDER

	elevatorOne.ID = 7
	elevatorOne.Direction = src.DIR_UP
	elevatorTwo.Floor = 0
	elevatorOne.InsideOrders[1] = ORDER	
	elevatorOne.OutsideOrders[1][0] = ORDER
	elevatorOne.OutsideOrders[2][1] = ORDER

	elevatorTwo.ID = 8
	elevatorTwo.Direction = src.DIR_UP
	elevatorTwo.Floor = 0
	elevatorTwo.InsideOrders[2] = ORDER
	elevatorTwo.OutsideOrders[3][0] = ORDER

	tools.PrintQueue(mergedQueue)

	queueArray := []src.ElevatorData{elevatorOne, elevatorTwo}

	assignedArray := assignOrders(queueArray, mergedQueue)

	println("After assignment: ")

	tools.PrintQueueArray(assignedArray)
}