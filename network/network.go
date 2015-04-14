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
var connectionStatus = make(map[string]bool) //map[IP]status

var timeoutLimit time.Duration = 2*time.Second
var sendPingInterval time.Duration = 150*time.Millisecond

type networkMessage struct {
	address string
	elevatorData []byte
}

// Channels made for communication across modules
var ChDataToQueue = make(chan src.ElevatorData)
var ChReadyToMerge = make(chan bool)
var ChQueueReadyToBeSent = make(chan src.ElevatorData)
var ChReceiptFromQueue = make(chan bool)

//TESTED:
func sendPing(broadcastConn *net.UDPConn){
	dataToSend := <- ChQueueReadyToBeSent		
	
	for {
		select {
			case outgoingData := <- ChQueueReadyToBeSent:
				dataToSend = outgoingData

			default:
				broadcastConn.Write(Pack(dataToSend))
				time.Sleep(sendPingInterval)
		}
	}
}

//TESTED:
func listenPing(chReceivedData chan networkMessage){
	
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
				
		chReceivedData <- networkMessage{ IPaddress, buffer[:lengthOfMessage] }
	}
}

// TESTED:
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

	return "192.168.0.102" // FOR WINDOWS AND TESTING ONLY!

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

// TESTED:
func timer(timeout chan bool) {
	for {
		time.Sleep(timeoutLimit)
		timeout <- true
	}
}


func NetworkHandler() {
	chTimeout := make(chan bool)
	chReceivedData := make(chan networkMessage)
	//chSendData := make(chan src.ElevatorData)

	IP = GetIPAddress()	
	broadcastConn := createBroadcastConn()

	go listenPing(chReceivedData)
	go sendPing(broadcastConn)
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

				connectionStatus[data.address] = true
				ChDataToQueue <- Unpack(data.elevatorData)
				<- ChReceiptFromQueue

			case <- chTimeout:
				//println("Elevators: ") // FOR TESTING!
				for address, status := range connectionStatus{
					println(address)
					switch status {
						case true:
							connectionStatus[address] = false
							 // FOR TESTING!
						case false:
							println("Is probably dead!")
							delete(connectionStatus, address)
					}
				}
			//	println(" ") // FOR TESTING
				ChReadyToMerge <- true

			default:
				time.Sleep(100*time.Millisecond)
		}
	}
}

// TODO:

// * implement sendPing() 									| OK
// * Can we send faster than TimeoutLimit?					| OK
// * send the received orders to queue 						| OK
// * make it so we don't receive our queue thru connection 	| OK 
// * send unique queues only to Queue modules				| OK


