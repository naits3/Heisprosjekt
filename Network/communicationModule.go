package network

import (
	"fmt"
	"net"
)

// needs to import meassageModule
// Com = Communication

func initNetwork(){
//	InitConnModule()
//	initComModule()
}

func initComModule(){
	// Initialiserer all channels(chReceive chan, chSend
	// Starts communication control
	// Starts also InitConnModule
}

func comHandler(){
	// Select which controls sending and receiving
	for{
		select{
			case: // Recive meassage 
				// Send pack for deserialization in meassageModule
			case: //offlinecase
		}
	}
}

func listenPack(port string, chReceive chan []byte){
	// Create a TCP listener and the specified port.
	// Accept connections and send the message to a buffer/channel
	// Repeat
}

func SendPack(pack []byte){
	//GetConnMap()
	//Loop through connMap and send pack to all elevators
}


