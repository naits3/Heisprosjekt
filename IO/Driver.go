package Driver

/*
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
*/
import "C"


//Define all neassacary enums

func InitializeDriver(){
	//start thread pollUserOrders
	//start thread pollFLoorSensors
	//start thread driverHandler
}


func driverHandler(){
	for{
		select{
			case order 	 := chUserOrder:
				//format order
				//send order
			case newFloor:= chArraivedFloor:
				//format newFloor
				//send newfloor 
			case commandfromControl := chCommandFromControl:
				//do what the command says
				//define a protocol when we implement the control module
		}
	}
}


func pollUserOrders(chUserOrder chan){
//Loops over all order getfunctions in the C kode
//Only button orders
}


func pollFLoorSensors(ch chanFloorSensor){
//Loops through all floor sensors and reads them

}

func formatOrder(order){
	//Format the order
}

func formatNewFloor(newFloor){
	//Format the newFLoor
}




