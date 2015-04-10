package controller

import (
	"testing"
	"Heisprosjekt/src"
	"time"
	//"Heisprosjekt/tools"
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
	testData1.OutsideOrders[1][src.BUTTON_DOWN] = src.ORDER
	testData1.OutsideOrders[1][src.BUTTON_UP] = src.ORDER
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
	// tools.PrintQueue(testData1)
	chNewOrdersFromQueue <- testData1

	<-chNewFloor

}

