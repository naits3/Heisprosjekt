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
	ElevatorData 			src.ElevatorData
}


func sendPing(broadcastConn *net.UDPConn, chOutgoingData chan src.ElevatorData){
	dataToSend := <- chOutgoingData		
	
	for {
		select {
			case outgoingData := <- chOutgoingData:
				dataToSend = outgoingData

			default:
				broadcastConn.Write(PackQueue(dataToSend))
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
	if err != nil {print("Error creating UDP") }
	return broadcastConn
}


func sendOrder(broadcastConn *net.UDPConn, order src.ButtonOrder) {
	broadcastConn.Write(PackOrder(order))
}


func GetIPAddress() string {
	addrs, err := net.InterfaceAddrs()
    if err != nil {
    	println(err)              
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
		time.Sleep(verifyConnectionInterval)
		timeout <- true
	}
}

func InitNetwork(	chIdentificationTQ chan string,
					chElevatorDataTQ chan ElevatorMessage,
					chElevatorDataFQ chan src.ElevatorData,
					chOrderTQ chan src.ButtonOrder,
					chOrderFQ chan src.ButtonOrder,
					chDisconElevatorTQ chan string) {

	chVerifyConnectedElevators := make(chan bool)
	chReceivedData := make(chan message)
	chSendData := make(chan src.ElevatorData)

	ourIP = GetIPAddress()	
	chIdentificationTQ <- ourIP
	broadcastConn := createBroadcastConn()
	go listenPing(chReceivedData)
	go sendPing(broadcastConn, chSendData)
	go timer(chVerifyConnectedElevators)

	outGoingData :=  <- chElevatorDataFQ
	chSendData <- outGoingData

	go networkManager(chIdentificationTQ, chElevatorDataTQ, chElevatorDataFQ, chOrderTQ, chOrderFQ, chDisconElevatorTQ, chVerifyConnectedElevators, chReceivedData, chSendData, broadcastConn)

}


func networkManager(chIdentificationTQ chan string,
					chElevatorDataTQ chan ElevatorMessage,
					chElevatorDataFQ chan src.ElevatorData,
					chOrderTQ chan src.ButtonOrder,
					chOrderFQ chan src.ButtonOrder,
					chDisconElevatorTQ chan string,
					chVerifyConnectedElevators chan bool,
					chReceivedData chan message,
					chSendData chan src.ElevatorData,
					broadcastConn *net.UDPConn) {

	for { 
		select {
			
			case order := <- chOrderFQ:
				sendOrder(broadcastConn, order)

			case receivedMessage := <- chReceivedData:
				if (receivedMessage.senderAddress == ourIP) {
					break
				}

				connectedElevators[receivedMessage.senderAddress] = true
				if (len(receivedMessage.data) > 50) {
					chElevatorDataTQ <- ElevatorMessage{receivedMessage.senderAddress, UnpackQueue(receivedMessage.data)}

				}else {
					chOrderTQ <- UnpackOrder(receivedMessage.data)
				}
				

			case outGoingData := <- chElevatorDataFQ:
				chSendData <- outGoingData

			

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