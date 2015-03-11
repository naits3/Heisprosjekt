package IO


import "/driver"

func InitializeDriver(chController chan int){
	//create all channels
	//start thread pollUserOrders
	//start thread pollFLoorSensors
	//start thread driverHandler

	chUserOrder := make(chan int)
	
	go driverHandler(chController)
	go pollUserOrders(chUserOrder)
}

func driverHandler(chController chan int){
	for{
		select{
			case order := chUserOrder:
				//format order
				//send order
				chController <- 1

			case newFloor:= chArraivedFloor:
				//format newFloor
				//send newfloor
				continue
			case commandfromControl := chCommandFromControl:
				//do what the command says
				//define a protocol when we implement the control module
				continue
		}
	}
}

func pollUserOrders(chUserOrder chan int){
//Loops over all order getfunctions in the C kode
//Only button orders

	i := 0
	for{
	 	i = 
	 	if i==1 {
	 		chUserOrder <-- 1
	 		i = 0
	 	}
	}
}


func pollFLoorSensors(chanFloorSensor chan int){
//Loops through all floor sensors and reads them
}

func formatOrder(order int) int{
	//Format the order
	return order
}

func formatNewFloor(newFloor int) int{
	//Format the newFLoor
	return newFloor
}




