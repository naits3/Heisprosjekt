package network

import (
	"testing"
	"time"
	"Heisprosjekt/src"
)

func TestGetIP(t *testing.T) {
	ip := GetIPAddress()
	println(ip)	
}


func TestBroadcastConn(t *testing.T){
	connection := createBroadcastConn()

	for {
		time.Sleep(time.Second)
		n, _ := connection.Write([]byte("IAM"))
		println(n)
	}
}


func TestListenPing(t *testing.T) {
	chNetworkMessage := make(chan message)
	
	go listenPing(chNetworkMessage)

	for {
		select {
			case data := <- chNetworkMessage:
				println("Got data from: ", data.senderAddress)

			default:
				time.Sleep(10*time.Millisecond)
		}
	}
}


func TestNetworkHandler(t *testing.T) {
	done := make(chan bool)
	go NetworkHandler()
	
	go func() {
		for {
			select {
				case <- ChElevatorDataToQueue:
					println("Got elevatorData")

				case order := <- ChOrderToQueue:
					println("t Q: -- floor ", order.Floor, ", type ", order.ButtonType)

				case <- ChReadyToMerge:
					println("Ready to merge. ConnectionStatus: ")
					for ID, status := range connectedElevators {
						println("ID: ", ID, "status: ", status)
					}
					println("")

				default:
					time.Sleep(100*time.Millisecond)
			}
		}

	}()

	var elevatorData src.ElevatorData
	elevatorData.Floor = 2
	elevatorData.Direction = src.DIR_DOWN
	
	ChQueueReadyToBeSent <- elevatorData

	//time.Sleep(2*time.Second)

	//ChOrderFromQueue <- src.ButtonOrder{2, src.BUTTON_DOWN}


	<- done
}