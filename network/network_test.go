package network

import (
	"testing"
	"time"
	//"Heisprosjekt/src"
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

func TestPacking(t *testing.T) {
	
}