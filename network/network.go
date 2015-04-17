package network

import (
	"net"
	"os"
	"strings"
	"time"
	"Heisprosjekt/src"
)

const PORT = "20019"	
var IP string
var connectedElevators = make(map[string]bool)

var timeoutLimit time.Duration = 1*time.Second
var sendMessageInterval time.Duration = 200*time.Millisecond

type message struct {
	senderAddress string
	data []byte
}


var ChElevatorDataToQueue = make(chan src.ElevatorData)
var ChReadyToMerge = make(chan bool)
var ChQueueReadyToBeSent = make(chan src.ElevatorData)


func sendPing(broadcastConn *net.UDPConn, chOutgoingData chan src.ElevatorData){
	dataToSend := <- chOutgoingData		
	
	for {
		select {
			case outgoingData := <- chOutgoingData:
				dataToSend = outgoingData

			default:
				broadcastConn.Write(Pack(dataToSend))
				time.Sleep(sendMessageInterval)
		}
	}
}


func listenPing(chReceivedMessage chan message){
	
	UDPAddr, _ := net.ResolveUDPAddr("udp",":"+PORT)
	var buffer []byte = make([]byte, 1024)
	conn, err := net.ListenUDP("udp", UDPAddr)

	if err != nil {
		println(err)
		return
	}

	defer conn.Close()
	for {
		lengthOfMessage, IPaddressAndPort, err := conn.ReadFromUDP(buffer)
		if err != nil {
			print(err)
			return
		}

		IPaddressAndPortArray := strings.Split(IPaddressAndPort.String(),":")
		IPaddress := IPaddressAndPortArray[0]
				
		chReceivedMessage <- message{ IPaddress, buffer[:lengthOfMessage] }
	}
}


func createBroadcastConn() *net.UDPConn{

	broadcastIP := GetIPAddress()

	ipArray := strings.Split(broadcastIP,".")
	ipArray[3] = "255"
	broadcastIP = strings.Join(ipArray,".")
	
	UDPAddr, err := net.ResolveUDPAddr("udp",broadcastIP + ":" + PORT)

	broadcastConn, err := net.DialUDP("udp",nil,UDPAddr)
	if err != nil {print("Error creating UDP") }// Error handling here}
	return broadcastConn
}


func GetIPAddress() string {

	addrs, err := net.InterfaceAddrs()
    if err != nil {
    	println(err)
		// error handle here                
    }

    for _, address := range addrs {
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
        	if ipnet.IP.To4() != nil {
            	return ipnet.IP.String()
        	}
    	}
    }
	
	println("Cannot resolve IP address! Exiting..")
	os.Exit(1)
	return ""
}

func timer(timeout chan bool) {
	for {
		time.Sleep(timeoutLimit)
		timeout <- true
	}
}


func NetworkHandler() {
	chVerifyConnectedElevators := make(chan bool)
	chReceivedData := make(chan message)
	chSendData := make(chan src.ElevatorData)

	IP = GetIPAddress()	
	broadcastConn := createBroadcastConn()

	go listenPing(chReceivedData)
	go sendPing(broadcastConn, chSendData)
	go timer(chVerifyConnectedElevators)

	for { 
		select {
			case receivedMessage := <- chReceivedData:
				
				alreadyReceived := false
				for storedAddress, status := range connectedElevators {
					if receivedMessage.senderAddress == storedAddress && status == true {
						alreadyReceived = true
					}
				}

				if (receivedMessage.senderAddress == IP || alreadyReceived) {
					break
				}

				connectedElevators[receivedMessage.senderAddress] = true
				ChElevatorDataToQueue <- Unpack(receivedMessage.data)

			case outGoingData := <- ChQueueReadyToBeSent:
				chSendData <- outGoingData

			case <- chVerifyConnectedElevators:
				for address, status := range connectedElevators{
					switch status {
						case true:
							connectedElevators[address] = false
						case false:
							delete(connectedElevators, address)
					}
				}
				ChReadyToMerge <- true
		}
	}
}

// TODO:

// Make an initNetwork
// Make the channels lokal (remove globals ... )