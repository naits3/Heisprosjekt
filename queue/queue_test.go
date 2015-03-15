package queue

import (
	"testing"
	"Heisprosjekt/src"
)


func TestMergeOrders(t *testing.T) {
	var testData1 src.ElevatorData
	var testData2 src.ElevatorData
	var testData3 src.ElevatorData
	//var mergedData elevatorData
	
	testData1.OutsideOrders[2][1] = ORDER
	testData2.OutsideOrders[2][1] = DELETE_ORDER
	testData3.OutsideOrders[2][1] = ORDER

	queueList := []src.ElevatorData{testData1, testData2, testData3}
	// ------- ONLY FOR PRINTING: -------------//
	println("Elevator 1:")
	src.PrintQueue(testData1)
	println("Elevator 2:")
	src.PrintQueue(testData2)
	println("Elevator 3:")
	src.PrintQueue(testData3)
	println(" ---------------- ")
	// ---------------------------------------//

	mergedData := mergeOrders(queueList)
	println("After Merge:")
	src.PrintQueue(mergedData)

}

func TestCalcTotalCost(t *testing.T) {
	var testData1 src.ElevatorData
	testData1.Direction = src.DIR_DOWN
	testData1.Floor = 1
	testData1.OutsideOrders[2][1] = ORDER
	testData1.OutsideOrders[1][1] = ORDER
	testData1.OutsideOrders[1][0] = ORDER
	//testData1.insideOrders[1] = ORDER

	src.PrintQueue(testData1)

	cost := calcTotalCost(testData1)
	println(cost)

	expect := 13
	if cost != expect {
		t.Error("expected", expect, "got", cost)
	} 

}