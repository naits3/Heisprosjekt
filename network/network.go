package network

import (
	"net"
	"strings"
	"time"
	"Heisprosjekt/src"
)

const PORT = "80"	
const IP = "192.168.0.102"
var connectionStatus = make(map[string]bool) //map[IP]status
var pingTimeLimit time.Duration = 10*time.Second

type networkMessage struct {
	address string
	elevatorData []byte
}

// Channels made for communication across modules
var ChDataToQueue = make(chan src.ElevatorData)
var ChReadyToMerge = make(chan bool)
var ChQueueReadyToBeSent = make(chan src.ElevatorData)

//TESTED:
func sendPing(broadcastConn *net.UDPConn, chSendData chan src.ElevatorData){
	for {
		select {
			case data := <- chSendData:
				broadcastConn.Write(Pack(data))
		}
	}
}

//TESTED:
func listenPing(chReceivedData chan networkMessage){
	UDPAddr, _ := net.ResolveUDPAddr("udp",":"+PORT)
	var buffer []byte = make([]byte, 1024)
	conn, err := net.ListenUDP("udp", UDPAddr)

	if err != nil {
			print(err)
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

				
		chReceivedData <- networkMessage{ IPaddress, buffer[:lengthOfMessage] }
	}
}

// TESTED:
func createBroadcastConn() *net.UDPConn{
	broadcastIP := GetIPAddress()
	
	switch(broadcastIP) {
		case "127.0.0.1":
			break
		default:
			ipArray := strings.Split(broadcastIP,".")
			ipArray[3] = "255"
			broadcastIP = strings.Join(ipArray,".")
	}

	UDPAddr, err := net.ResolveUDPAddr("udp",broadcastIP + ":" + PORT)

	broadcastConn, err := net.DialUDP("udp",nil,UDPAddr)
	if err != nil {print("Error creating UDP") }// Error handling here}
	return broadcastConn
}

// TESTED:
func GetIPAddress() string{
	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
	    addrs, _ := i.Addrs()
	    // handle err
	    for _, addr := range addrs {
	        switch addressType:= addr.(type) {
	        case *net.IPAddr:
	            if(addressType.IP.String() != "0.0.0.0"){
	            	return addressType.IP.String()
	            	
	        	}
	        }
	    }
	}
	return "127.0.0.1"
	// Vi får et problem dersom heisen vår er frakoblet i oppstart, og skal koble seg på nettet seinere.
	// Én løsning er å slette denne funksjonen, og sette IP direkte. 
	// ... Eller ikke. Hvordan vil vi da kunne sende kø i offline?
}

// TESTED:
func timer(timeout chan bool) {
	for {
		time.Sleep(pingTimeLimit)
		timeout <- true
	}
}


func NetworkHandler() {
	chTimeout := make(chan bool)
	chReceivedData := make(chan networkMessage)
	chSendData := make(chan src.ElevatorData)

	broadcastConn := createBroadcastConn()

	go listenPing(chReceivedData)
	go sendPing(broadcastConn, chSendData)
	go timer(chTimeout)

	for { 
		select {
			case data := <- chReceivedData:
				alreadyReceived := false
				for storedAddress, status := range connectionStatus {
					if data.address == storedAddress && status == true {
						alreadyReceived = true
					}
				}

				if (data.address == IP || alreadyReceived) {
					break
				}

				println("received from: ", data.address)
				connectionStatus[data.address] = true
				ChDataToQueue <- Unpack(data.elevatorData)

			case outGoingData := <- ChQueueReadyToBeSent:
				chSendData <- outGoingData

			case <- chTimeout:
				for address, status := range connectionStatus{
					println(address)
					switch status {
						case true:
							connectionStatus[address] = false
						case false:
							delete(connectionStatus, address)
					}
				}

				ChReadyToMerge <- true
		}
	}
}

// TODO:

// * implement sendPing() 									| OK
// * handle connections when a elevator does not respond	|
// * send the received orders to queue 						| OK
// * make it so we don't receive our queue thru connection 	| 
// * set IPadress manually									|
// * send unique queues only to Queue modules				| 