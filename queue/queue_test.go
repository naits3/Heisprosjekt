package queue

import (
	"testing"
	"time"
	"Heisprosjekt/src"
	"Heisprosjekt/tools"
)

func TestCalcOrderCost(t *testing.T) {
	var elevator src.ElevatorData
	elevator.Floor = 0
	order := src.ButtonOrder{1, src.BUTTON_DOWN}
	println(calcOrderCost(elevator, order))
}

func TestAssignOrder(t *testing.T) {
	ourID = "192.168.0.101"
	var elevator1 src.ElevatorData
	var elevator2 src.ElevatorData
	elevatorQueues["192.168.0.101"] = elevator1
	
	tmp := elevatorQueues["192.168.0.101"]
	tmp.Floor = 2
	elevatorQueues["192.168.0.101"] = tmp

	elevatorQueues["192.168.0.102"] = elevator2

	order := src.ButtonOrder{3, src.BUTTON_UP}

	assignOrder(elevatorQueues, order)

	tools.PrintQueue(elevatorQueues[ourID])


	
}

func TestMergeOrders(t *testing.T) {
	var elevator1 src.ElevatorData
	var elevator2 src.ElevatorData
	var elevator3 src.ElevatorData

	ID1 := "192.168.0.101"
	ID2 := "192.168.0.102"
	ID3 := "192.168.0.103"

	elevatorQueues[ID1] = elevator1
	elevatorQueues[ID2] = elevator2
	elevatorQueues[ID3] = elevator3

	addOrder(ID1, src.ButtonOrder{0, src.BUTTON_UP})
	addOrder(ID2, src.ButtonOrder{2, src.BUTTON_UP})
	addOrder(ID3, src.ButtonOrder{2, src.BUTTON_DOWN})
	addOrder(ID3, src.ButtonOrder{3, src.BUTTON_INSIDE})

	tools.PrintQueue(mergeOrders(elevatorQueues))

	deleteOrder(ID1, src.ButtonOrder{0, src.BUTTON_UP})

	tools.PrintQueue(mergeOrders(elevatorQueues))
}

func TestQueueManager(t *testing.T) {
	chFloorFromController 		:= make(chan int)
	chOrderFromController		:= make(chan src.ButtonOrder)
	chDirectionFromController	:= make(chan int)
	chFinishedFromController	:= make(chan bool)
	chGlobalOrdersToController	:= make(chan src.ElevatorData)
	chNewNextFloorFromQueue		:= make(chan int)
	done 						:= make(chan bool)



	go queueManager(	chFloorFromController, 
				chOrderFromController, 
				chDirectionFromController, 
				chFinishedFromController, 
				chGlobalOrdersToController,
				chNewNextFloorFromQueue)

	go func(){
		for {
			select {
			
				case globOrd := <- chGlobalOrdersToController:
					println(" --------- ORDERS ---------- ")
					println("Global Orders: ")
					tools.PrintQueue(globOrd)
					
					for id, queue := range elevatorQueues {
						println(id)
						tools.PrintQueue(queue)
					}

				case floor := <- chNewNextFloorFromQueue:
					println("next floor: ", floor)

				default:
					time.Sleep(100*time.Millisecond)
			}
		}
	}()

	chFloorFromController <- 3

	time.Sleep(3*time.Second)
	chOrderFromController <- src.ButtonOrder{0, src.BUTTON_UP}

	//time.Sleep(2*time.Second)
	//chFinishedFromController <- true

	<- done
}