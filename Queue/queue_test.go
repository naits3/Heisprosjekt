package Queue

import (
	"testing"
)

// OBS! For aa faa print til aa fungere, maa funksjonsnavnene vaere global! Stor capital letter

func printQueue(queueData elevatorData) {
	for row := 0; row < FLOORS; row ++ {
		for col := 0; col < 2; col ++ {
			print(" ",queueData.outsideOrders[FLOORS-1-row][col])
		}
		println(" ")
	}
	println(" ")
}


func TestMergeOrders(t *testing.T) {
	var testData1 elevatorData
	var testData2 elevatorData 
	var testData3 elevatorData
	//var mergedData elevatorData
	
	testData1.outsideOrders[2][1] = ORDER
	testData2.outsideOrders[2][1] = DELETE_ORDER
	testData3.outsideOrders[2][1] = ORDER

	queueList := []elevatorData{testData1, testData2, testData3}
	// ------- ONLY FOR PRINTING: -------------//
	println("Elevator 1:")
	printQueue(testData1)
	println("Elevator 2:")
	printQueue(testData2)
	println("Elevator 3:")
	printQueue(testData3)
	println(" ---------------- ")
	// ---------------------------------------//

	mergedData := MergeOrders(queueList)
	println("After Merge:")
	printQueue(mergedData)

}