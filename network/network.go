package network

import (
	"net"
	"os"
	"strings"
	"time"
	"Heisprosjekt/src"
)

const PORT = "20013"	
var ourIP string
var connectedElevators = make(map[string]bool)

var verifyConnectionInterval time.Duration = 2*time.Second
var sendMessageInterval time.Duration = 80*time.Millisecond

type message struct {
	senderAddress 	string
	data 			[]byte
}

type ElevatorMessage struct {
	SenderAddress 	string
	ElevatorData 	src.ElevatorData
}


func sendElevatorData(broadcastConn *net.UDPConn, chSendElevatorData chan src.ElevatorData){
	dataToSend := <- chSendElevatorData		
	
	for {
		select {
			case elevatorData := <- chSendElevatorData:
				dataToSend = elevatorData

			default:
				broadcastConn.Write(PackElevatorData(dataToSend))
				time.Sleep(sendMessageInterval)
		}
	}
}


func listenForMessage(chReceivedMessage chan message){
	
	UDPAddr, _ := net.ResolveUDPAddr("udp",":"+PORT)
	var buffer []byte = make([]byte, 1024)
	conn, err := net.ListenUDP("udp", UDPAddr)

	if err != nil {
		println("Cannot create a connection! Exiting..")
		os.Exit(1)
	}

	defer conn.Close()
	for {
		lengthOfMessage, IPaddressAndPort, _ := conn.ReadFromUDP(buffer)
		
		IPaddressAndPortArray := strings.Split(IPaddressAndPort.String(),":")
		IPaddress := IPaddressAndPortArray[0]
		
		chReceivedMessage <- message{ IPaddress, buffer[:lengthOfMessage] }
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
	packedOrder := PackOrder(order)
	broadcastConn.Write(packedOrder)
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
	chReceivedMessage 			:= make(chan message)
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
					chReceivedMessage chan message,
					chSendElevatorData chan src.ElevatorData,
					broadcastConn *net.UDPConn) {

	for { 
		select {
			
			case order := <- chOrderFQ:
				sendOrder(broadcastConn, order)

			case receivedMessage := <- chReceivedMessage:
				if (receivedMessage.senderAddress == ourIP) {
					break
				}

				connectedElevators[receivedMessage.senderAddress] = true
			
				if (len(receivedMessage.data) > 50) {
					chElevatorDataTQ <- ElevatorMessage{receivedMessage.senderAddress, UnpackElevatorData(receivedMessage.data)}

				}else {
					chOrderTQ <- UnpackOrder(receivedMessage.data)
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