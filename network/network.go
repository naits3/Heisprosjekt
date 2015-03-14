package network

import (
	"net"
	"strings"
	"time"
	"Heisprosjekt/src"
)

var connStorage map[string]net.Conn = nil
const PORT = "80"	


var connectionStatus = make(map[string]bool) //map[IP]status
var pingTimeLimit time.Duration = time.Second

// Channels made for communication across modules
var ChDataToQueue = make(chan src.ElevatorData)
var ChReadyToMerge = make(chan bool)
var ChQueueReadyToBeSent = make(chan src.ElevatorData)

//TESTED:
func sendPing(broadcastConn *net.UDPConn){
	for {
		select {
			case data := <- ChQueueReadyToBeSent:
				broadcastConn.Write(Pack(data))
		}
	}
}

//TESTED:
func listenPing(chReceivedData chan []byte, chReceivedIPaddress chan string){
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
		
		chReceivedIPaddress <- IPaddress
		chReceivedData <- buffer[:lengthOfMessage]
	}
}

// TESTED:
func createBroadcastConn() *net.UDPConn{
	broadcastIP := getBroadcastIP()
	UDPAddr, err := net.ResolveUDPAddr("udp",broadcastIP + ":" + PORT)

	broadcastConn, err := net.DialUDP("udp",nil,UDPAddr)
	if err != nil {print("Error creating UDP") }// Error handling here}
	return broadcastConn
}

// TESTED: 
func getBroadcastIP() string{
	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
	    addrs, _ := i.Addrs()
	    // handle err
	    for _, addr := range addrs {
	        switch addressType:= addr.(type) {
	        case *net.IPAddr:
	            if(addressType.IP.String() != "0.0.0.0"){
	            	localAddr := addressType.IP.String()
	            	localAddrSplitted := strings.Split(localAddr,".")
	            	localAddrSplitted[3] = "255"
	            	broadcastAddr := strings.Join(localAddrSplitted,".")
	            	return(broadcastAddr)
	        	}
	        }
	    }
	}
	return "is_offline" // EDIT THIS?
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
	chReceivedData := make(chan []byte)
	chReceivedIPAddress := make(chan string)

	go listenPing(chReceivedData, chReceivedIPAddress)
	go timer(chTimeout)

	for {
		select {
			case data := <- chReceivedData:
				ChDataToQueue <- Unpack(data)

			case IPAddress := <- chReceivedIPAddress:
				connectionStatus[IPAddress] = true

			case <- chTimeout:
				for address, status := range connectionStatus{
					//println(address) // FOR TESTING ONLY
					switch status {
						case true:
							connectionStatus[address] = false
						case false:
							delete(connectionStatus, address)
					}
				}
				// Send an OK-signal to the queue here when connections are handled
				ChReadyToMerge <- true
		}
	}
}

// TODO:

// * implement sendPing() 									| OK
// * handle connections when a elevator does not respond	|
// * send the received orders to queue 						| OK