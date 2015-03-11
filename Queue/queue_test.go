package Queue

import (
	"testing"
)

// OBS! For aa faa print til aa fungere, maa funksjonsnavnene vaere global! Stor capital letter

func printQueue(queueData elevatorData) {
	for row := 0; row < FLOORS; row ++ {
		for col := 0; col < 2; col ++ {
			print(queueData.outsideOrders[FLOORS-1-row][col])
		}
		println(" ")
	}
}


func TestMergeOrders(t *testing.T) {
	var testData1 elevatorData
	//var testData2 elevatorData 
	testData1.outsideOrders[0][0] = ORDER

	// ------- ONLY FOR PRINTING: -------------//
	printQueue(testData1)
	// ---------------------------------------//

}