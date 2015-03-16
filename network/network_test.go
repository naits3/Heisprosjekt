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


func TestSendPing(t *testing.T){
	bConn := createBroadcastConn()
	var testData src.ElevatorData

	go sendPing(bConn)
	for {
		time.Sleep(time.Second)
		ChQueueReadyToBeSent <- testData
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