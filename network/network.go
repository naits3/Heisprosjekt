package network

import (
	"net"
	"os"
	"strings"
	"time"
	"Heisprosjekt/src"
)

/* Name-conventions:
FQ = From Queue
TQ = To Queue
*/

const PORT = "20013"	
var ourIP string
var connectedElevators = make(map[string]bool)

type Message struct {
	SenderAddress 	string
	MessageType		string
	ElevatorData 	src.ElevatorData
	Order 			src.ButtonOrder
}

type ElevatorMessage struct {
	SenderAddress 	string
	ElevatorData 	src.ElevatorData
}


func sendElevatorData(broadcastConn *net.UDPConn, chSendElevatorData chan src.ElevatorData){
	sendMessageInterval := 80*time.Millisecond
	dataToSend := <- chSendElevatorData		
	
	for {
		select {
			case elevatorData := <- chSendElevatorData:
				dataToSend = elevatorData

			default:
				msg := Message{ourIP, "elevatorData", dataToSend, src.ButtonOrder{}}
				packedMsg := PackMessage(msg)
				broadcastConn.Write(packedMsg)
				time.Sleep(sendMessageInterval)
		}
	}
}


func listenForMessage(chReceivedMessage chan Message){
	
	UDPAddr, _ := net.ResolveUDPAddr("udp",":"+PORT)
	var buffer []byte = make([]byte, 1024)
	conn, err := net.ListenUDP("udp", UDPAddr)

	if err != nil {
		println("Cannot create a connection! Exiting..")
		os.Exit(1)
	}

	defer conn.Close()
	for {
		lengthOfMessage, _, _ := conn.ReadFromUDP(buffer)

		msg := UnpackMessage(buffer[:lengthOfMessage])
		chReceivedMessage <- msg
	}
}


func createBroadcastConn() *net.UDPConn{

	ipAddress := GetIPAddress()

	ipArray := strings.Split(ipAddress,".")
	ipArray[3] = "255"
	broadcastIP := strings.Join(ipArray,".")
	
	UDPAddr, _ := net.ResolveUDPAddr("udp",broadcastIP + ":" + PORT)

	broadcastConn, err := net.DialUDP("udp",nil,UDPAddr)
	if err != nil {
		println("Cannot create a connection! Exiting..")
		os.Exit(1)
	}
	return broadcastConn
}


func sendOrder(broadcastConn *net.UDPConn, order src.ButtonOrder) {
	msg := Message{ourIP, "order", src.ElevatorData{}, order}
	packedMsg := PackMessage(msg)
	broadcastConn.Write(packedMsg)
}


func GetIPAddress() string {
	addrs, _ := net.InterfaceAddrs()
    
    for _, address := range addrs {
    	ipnet, ok := address.(*net.IPNet)
        if (ok && !ipnet.IP.IsLoopback()) {
        	if (ipnet.IP.To4() != nil) {
            	return ipnet.IP.String()
        	}
    	}
    }
	
	println("Cannot resolve IP address! Exiting..")
	os.Exit(1)
	return ""
}


func verifyConnectionTimer(chVerifyConnectedElevators chan bool) {
	verifyConnectionInterval := 2*time.Second

	for {
		time.Sleep(verifyConnectionInterval)
		chVerifyConnectedElevators <- true
	}
}

func InitNetwork(	chIdentificationTQ chan string,
					chElevatorDataTQ chan ElevatorMessage,
					chElevatorDataFQ chan src.ElevatorData,
					chOrderTQ chan src.ButtonOrder,
					chOrderFQ chan src.ButtonOrder,
					chDisconElevatorTQ chan string) {


	chVerifyConnectedElevators 	:= make(chan bool)
	chReceivedMessage 			:= make(chan Message)
	chSendElevatorData 			:= make(chan src.ElevatorData)

	ourIP = GetIPAddress()	
	chIdentificationTQ <- ourIP
	broadcastConn := createBroadcastConn()
	go listenForMessage(chReceivedMessage)
	go sendElevatorData(broadcastConn, chSendElevatorData)
	go verifyConnectionTimer(chVerifyConnectedElevators)

	elevatorData :=  <- chElevatorDataFQ
	chSendElevatorData <- elevatorData

	go networkManager(chIdentificationTQ, chElevatorDataTQ, chElevatorDataFQ, chOrderTQ, chOrderFQ, chDisconElevatorTQ, chVerifyConnectedElevators, chReceivedMessage, chSendElevatorData, broadcastConn)

}


func networkManager(chIdentificationTQ chan string,
					chElevatorDataTQ chan ElevatorMessage,
					chElevatorDataFQ chan src.ElevatorData,
					chOrderTQ chan src.ButtonOrder,
					chOrderFQ chan src.ButtonOrder,
					chDisconElevatorTQ chan string,
					chVerifyConnectedElevators chan bool,
					chReceivedMessage chan Message,
					chSendElevatorData chan src.ElevatorData,
					broadcastConn *net.UDPConn) {

	for { 
		select {
			
			case order := <- chOrderFQ:
				sendOrder(broadcastConn, order)

			case receivedMessage := <- chReceivedMessage:
				if (receivedMessage.SenderAddress == ourIP) {
					break
				}

				connectedElevators[receivedMessage.SenderAddress] = true
				switch receivedMessage.MessageType {
					case "elevatorData":
						chElevatorDataTQ <- ElevatorMessage{receivedMessage.SenderAddress, receivedMessage.ElevatorData}

					case "order":
						chOrderTQ <- receivedMessage.Order
				}

			case elevatorData := <- chElevatorDataFQ:
				chSendElevatorData <- elevatorData		

			case <- chVerifyConnectedElevators:
				for address, status := range connectedElevators{
					switch status {
						case true:
							connectedElevators[address] = false
						case false:
							delete(connectedElevators, address)
							chDisconElevatorTQ <- address

					}
				}
		}
	}
}