package network

import (
	"net"
	"os"
	"strings"
	"time"
	"Heisprosjekt/src"
)

const PORT = "20013"	
var IP string
var connectedElevators = make(map[string]bool)

var timeoutLimit time.Duration = 2*time.Second
var sendMessageInterval time.Duration = 200*time.Millisecond

type message struct {
	senderAddress 	string
	messageType 	string
	data 			[]byte
}

type QueueMessage struct {
	SenderAddress 	string
	Data 			src.ElevatorData
}


var ChElevatorDataToQueue 	= make(chan QueueMessage)
var ChOrderToQueue 			= make(chan src.ButtonOrder)
var ChIDFromNetwork			= make(chan string)
var ChQueueReadyToBeSent 	= make(chan src.ElevatorData, 2)
var ChOrderFromQueue		= make(chan src.ButtonOrder, 2)
var ChLostElevator			= make(chan string)


func sendPing(broadcastConn *net.UDPConn, chOutgoingData chan src.ElevatorData){
	dataToSend := <- chOutgoingData
	
	for {
		select {
			case outgoingData := <- chOutgoingData:
				dataToSend = outgoingData

			default:
				UnpackedMessage := message{IP, "elevatorData", PackQueue(dataToSend)}
				packedMessage := PackMessage(UnpackedMessage)
				broadcastConn.Write(packedMessage)
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
		//lengthOfMessage, IPaddressAndPort, err := conn.ReadFromUDP(buffer)
		lengthOfMessage, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			print(err)
			return
		}

		//IPaddressAndPortArray := strings.Split(IPaddressAndPort.String(),":")
		//IPaddress := IPaddressAndPortArray[0]
		
		//chReceivedMessage <- message{ IPaddress, buffer[:lengthOfMessage] }
		chReceivedMessage <- UnpackMessage(buffer[:lengthOfMessage])
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


func sendOrder(broadcastConn *net.UDPConn, message []byte) {
	broadcastConn.Write(message)
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
	ChIDFromNetwork <- IP
	broadcastConn := createBroadcastConn()

	go listenPing(chReceivedData)
	go sendPing(broadcastConn, chSendData)
	go timer(chVerifyConnectedElevators)

	outGoingData :=  <- ChQueueReadyToBeSent
	chSendData <- outGoingData

	for { 
		select {
			
			case receivedMessage := <- chReceivedData:
				if (receivedMessage.senderAddress == IP) {
					break
				}

				connectedElevators[receivedMessage.senderAddress] = true
				switch receivedMessage.messageType {
					
					case "elevatorData":
						ChElevatorDataToQueue <- QueueMessage{receivedMessage.senderAddress, UnpackQueue(receivedMessage.data)}

					case "order":
						ChOrderToQueue <- UnpackOrder(receivedMessage.data)
				}

			case outGoingData := <- ChQueueReadyToBeSent:
				chSendData <- outGoingData

			case order := <- ChOrderFromQueue:
				message := message{IP, "order", PackOrder(order)}
				packedMessage := PackMessage(message)
				sendOrder(broadcastConn, packedMessage)

			case <- chVerifyConnectedElevators:
				for address, status := range connectedElevators{
					switch status {
						case true:
							connectedElevators[address] = false
						case false:
							delete(connectedElevators, address)
							ChLostElevator <- address

					}
				}
		}
	}
}



// TODO:

// * Slå seg av dersom den ikke får laget connections?

// * Lag en generell meldingsstruct. 