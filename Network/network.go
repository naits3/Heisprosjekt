package Network

import (
	"net"
	"strings"
	//"time"
)

const FLOORS = 4

var connStorage map[string]net.Conn = nil
const PORT = "3000"	

type elevatorData struct { 
	IP 				int
	floor 			int
	direction 		int
	outsideOrders [FLOORS][2]int
	insideOrders  [FLOORS]int
}


func sendPing(udpConn net.Conn){
	
}


//TESTED:
// - need to test with multiple clients
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