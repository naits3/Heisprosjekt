package controller

import (
	"testing"
	"Heisprosjekt/src"
	"time"
	"Heisprosjekt/tools"
)

func TestInitController(t *testing.T){
	
	InitController()

	u := 1
	if u != 1{
		t.Errorf("Error have happend")
	}
}


func TestLSetLights(t *testing.T){
	
	InitController()
	time.Sleep(1*time.Second)

	var testData1 src.ElevatorData
	testData1.OutsideOrders[2][src.BUTTON_DOWN] = src.ORDER
	testData1.OutsideOrders[0][src.BUTTON_UP] = src.ORDER
	testData1.OutsideOrders[3][src.BUTTON_DOWN] = src.ORDER
	testData1.ID = 1

	setLights(testData1)

	u := 1
	if u != 1{
		t.Errorf("Error have happend")
	}
}

func TestSetLights(t *testing.T){
	InitController()

	var testData1 src.ElevatorData
	testData1.OutsideOrders[2][src.BUTTON_UP] = src.ORDER
	testData1.OutsideOrders[1][src.BUTTON_DOWN] = src.ORDER
	testData1.OutsideOrders[1][src.BUTTON_UP] = src.ORDER
	testData1.ID = 1



	time.Sleep(1*time.Second)
	 tools.PrintQueue(testData1)
	chNewOrdersFromQueue <- testData1

	time.Sleep(1*time.Second)
}

func TestReciveOrder(t *testing.T){
	InitController()
	
	for{
		println("Push a button:")
		order := <- chNewOrder
		println(order.ButtonType)
		println(order.Floor)
		break
 	}
}

func TestNextFloor(t *testing.T) {
	InitController()
	println("waiting for launch")
	println()
	time.Sleep(time.Second)

	// Send in your floor here:

	chNewNextFloorFromQueue <- 1

	<- chNewFloor

}

