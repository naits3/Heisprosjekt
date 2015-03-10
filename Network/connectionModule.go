package main

import (
	"net"
	"strings"
	"time"
	//"fmt"
)

var connStorage map[string]net.Conn = nil
var port string	= "3000"						//Decide a listening port!!

// conn is connection

func InitConnModule(){
	//init all channels
	//Create broadcast connection
	//Thread connDurationHandler()
	//Thread listenPing()
	// Thread sendPing()
}

func addConn(ip string){
	//adds a connection to map with ip as key
}

func removeConn(ip string){
	//remove conn from map with ip string
}

func GetConnMap() {
	// -- RETURNS map[string]net.Conn --
	// Sender alle TCP til kommunikasjonskontroll
	// Må ha variable kontroll med addConn.. og deleteConn..
}

func sendPing(udpConn net.Conn){
	message := "I AM ALIVE"
	timeDur := 500*time.Millisecond
	
	for _ = range time.Tick(timeDur){
		udpConn.Write([]byte(message)) // Send JSON-melding her
	}
}

func listenPing(port string, chListenPing chan string, quit chan bool) string{
	addr, err := net.ResolveUDPAddr("udp","3000")
	if err != nil{return "0"} // Handle error
	conn, err := net.ListenUDP("udp", addr)
	if err != nil{return "0"} // Handle error

	var buffer []byte = make([]byte, 1500)
	
	for{
		select {
		case <- quit:
			conn.Close()
			return "hade"
		default:
			_, senderAddr, err := conn.ReadFromUDP(buffer)
			if err != nil{return "0"} // Handle error
			chListenPing <- senderAddr.String()
		}
	}

	// listning for broadcast signals
	// if valid signal send to connDurationHandler
}

func connDurationHandler(){
	// start new timer for 1 sec
	// Get all keys
	// for{
	// 	select{
	// 	case time:=<-Timer:
	// 		//check if the list are empty, if not remove ip from map
	// 		//start new timer for 1 sec
	// 	//case: //gets ip string from listenPing() channel
	// 		//check if it is in map
	// 		//if in map remove from list
	// 		//if not in map add new connection to map

	// 	}
	//}
}

func createBroadcastConn() net.Conn{
	broadcastIP := getBroadcastIP()
	
	// println(broadcastIP)
	conn, err := net.Dial("udp",broadcastIP+":"+port)
	if err != nil {print("Error creating UDP") }// Error handling here}
	return conn

}

func GetPort() string{
	return port
}

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


func main(){
	chListenPing := make(chan string) // Buffersize?
	chQuit := make(chan bool)

	go listenPing(port, chListenPing, chQuit)

	for {
		select {
			case elevatorAddr := <- chListenPing:
				println(elevatorAddr)

			case <- time.After(time.Second):
				println("Time's up moddafucka")
		}
	}	
}