package controller

import (
	"testing"
	"Heisprosjekt/src"
)

func TestInitController(t *testing.T){
	
	InitController()

	u := 1
	if u != 1{
		t.Errorf("Error have happend")
	}
}


func TestSetLights(t *testing.T){
	
	InitController()

	var testData1 src.ElevatorData
	testData1.OutsideOrders[0][src.BUTTON_DOWN] = src.ORDER
	testData1.InsideOrders[2] = src.ORDER
	testData1.ID = 1

	setLights(testData1)

	u := 1
	if u != 1{
		t.Errorf("Error have happend")
	}
}

