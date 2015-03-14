package network

import (
	"testing"
	"time"
	"Heisprosjekt/src"
)

func TestGetIP(t *testing.T) {
	ip := getBroadcastIP()
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
	chPing := make(chan []byte)
	chAddress := make(chan string)
	
	go listenPing(chPing, chAddress)

	for {
		select {
			case data := <- chPing:
				println(string(data))
			case addr := <- chAddress:
				println("address: ",addr)
		}
	}
}

func send(){
	connection := createBroadcastConn()

	for {
		time.Sleep(time.Second)
		connection.Write([]byte("IAM"))
		
	}
}

func TestNetworkHandler(t *testing.T) {
	NetworkHandler()
	

	// test := make(map[string]bool)
	// test["hei"] = true
	// //delete(test,"hei")

	// println("1:")
	// for key, value := range test{
		
	// 	println(key+":", value)
	// 	delete(test, key)
	// }

	// println("2:")
	// for key, value := range test{
	// 	println(key+":", value)
	// 	//delete(test, key)
	// }
}

func TestSendAndReceive(t *testing.T) {
	// SENDING:
	

	// LISTENING:
	chPing := make(chan []byte)
	chAddress := make(chan string)
	
	//go send()
	go listenPing(chPing, chAddress)

	for {
		select {
			case data := <- chPing:
				println(string(data))
			case addr := <- chAddress:
				println("address: ",addr)
		}
	}

	// Conclusion: We listen to ourself if we broadcast and listen
}

func TestSendPing(t *testing.T){
	bConn := createBroadcastConn()
	var testData src.ElevatorData

	go sendPing(bConn)
	for {
		time.Sleep(time.Second)
		ChQueueReadyToBeSent <- testData
	}
}