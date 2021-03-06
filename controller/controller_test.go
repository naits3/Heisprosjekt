package controller

import (
	"testing"
	//"Heisprosjekt/src"
	//"time"
	//"Heisprosjekt/tools"
)

// func TestInitController(t *testing.T){
	
// 	InitController()

// 	u := 1
// 	if u != 1{
// 		t.Errorf("Error have happend")
// 	}
// }


// func TestLSetLights(t *testing.T){
	
// 	InitController()
// 	time.Sleep(1*time.Second)

// 	var testData1 src.ElevatorData
// 	testData1.OutsideOrders[2][src.BUTTON_DOWN] = src.ORDER
// 	testData1.OutsideOrders[0][src.BUTTON_UP] = src.ORDER
// 	testData1.OutsideOrders[3][src.BUTTON_DOWN] = src.ORDER
// 	testData1.ID = 1

// 	setLights(testData1)

// 	u := 1
// 	if u != 1{
// 		t.Errorf("Error have happend")
// 	}
// }

// func TestSetLights(t *testing.T){
// 	InitController()

// 	var testData1 src.ElevatorData
// 	testData1.OutsideOrders[2][src.BUTTON_UP] = src.ORDER
// 	testData1.OutsideOrders[1][src.BUTTON_DOWN] = src.ORDER
// 	testData1.OutsideOrders[1][src.BUTTON_UP] = src.ORDER
// 	testData1.ID = 1

// 	time.Sleep(1*time.Second)
// 	tools.PrintQueue(testData1)
// 	chNewOrdersFromQueue <- testData1
// 	time.Sleep(1*time.Second)
// }

// func TestReciveOrder(t *testing.T){
// 	InitController()
	
// 	for{
// 		println("Push a button:")
// 		order := <- chNewOrder
// 		println(order.ButtonType)
// 		println(order.Floor)
//  	}
// }

// func TestNextFloor(t *testing.T) {
	
// 	println("waiting for launch")
// 	// Send in your floor here:

	
// 	<- chOrderIsFinished
// 	chNewNextFloorFromQueue <- 2
// 	<- chOrderIsFinished
// 	chNewNextFloorFromQueue <- 1
// 	<- chOrderIsFinished
// 	chNewNextFloorFromQueue <- 0
// 	<- chOrderIsFinished
 
// }

// func TestFindElevatorDirection(t *testing.T){
// 	currentFloor = 2
// 	nextFloor 	= 2

// 	println(findElevatorDirection())
// }

func TestControllerFunctions(t *testing.T){

	stop := make(chan bool)

	InitController()

	//chDestinationFloorFromQueue <- 3

	go func(){
		for {
			select {
			case floor := <- chFloorToQueue:
				println("Current floor: ",floor)

			case order := <- chOrderToQueue:
				println("Order floor: ",order.Floor)
				println("Order ButtonType: ",order.ButtonType)
			 	chDestinationFloorFromQueue <- order.Floor
			case direction := <- chDirectionToQueue:
				println("Direction: ", direction)
			
			case finished := <- chOrderFinishedToQueue:
				println("Finised :", finished)
			
			}
		}
	}()

	<- stop
}