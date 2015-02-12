
package network

import (
	"fmt"
	"net"
)

var connectionMap map[string]net.Conn = nil

func initNetwork(){
	// legger til broadcast i map med key BROADCAST
	// legger til andre connections med funksjonen requestConnections
	//go ping()
	//go listenForConnections()
	//go requestConnections()
}

func GetHostIP() string{
	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
	    addrs, _ := i.Addrs()
	    // handle err
	    for _, addr := range addrs {
	        switch v:= addr.(type) {
	        case *net.IPAddr:
	            if(v.IP.String() != "0.0.0.0"){
	            	return(v.IP.String())
	        	}
	        }
	    }
	}
	return "is_offline"
}

func SendPackToAll(){
	//Sender til alle TCP adresser i dictonary
}

func receivePack(port string, receive chan []byte) []byte {
	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		return []byte("0")
	}

	for {
		var buffer []byte = make([]byte, 1024)
		conn, err := listener.Accept()
		packetSize, err := conn.Read(buffer)

		if err != nil {
			return []byte("0")
		}

		if packetSize > 0{
			receive <- buffer[0:packetSize]
		}
	}
}

func sendPack(pack []byte, host string, chSend chan bool){
	addr, _ := net.ResolveTCPAddr("tcp",host)
	conn, err := net.DialTCP("tcp", nil, addr)

	if err != nil{
		return
	}
	
	_, err = conn.Write(pack)
	
	if err != nil {
		return
	}

	chSend <- true
}


func listenForConnections(port int){
	addr := net.UDPAddr{
		Port: port,
	}

	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		fmt.Println("Error listening to UDP: ",err)
		return
	}

	var buffer []byte = make([]byte, 1024)
	
	defer conn.Close()
	for {
		n, address, err := conn.ReadFromUDP(buffer)

		if err != nil {
			fmt.Println("Error reading from UDP: ",err)
			return
		}

		if address != nil && string(buffer[0:n])=="CONNECTME" {
			// SEND OK HERE
		}
	}
}

func requestConnection(ip string, port string, ){
	//UDP broadcast en gang ved oppstart
	//hvis ikke connection prøv igjen altså hvis det kommer en error
}

func ping(){
	//UDP broadcast
}

func addConnection(ip string){
	//Legger til ip som hash variabel med en connection TCP
}

func deleteConnection(ip string){
	//Fjerner en connection fra dictonary for TCP
}


func main(){
	port := "20020"
	host := "78.91.38.8"+":"+port

	chReceive := make(chan []byte)
	chSend := make(chan bool)	
	
	//packetBuffer := make([]byte, 1024)

	go ReceivePack(port, chReceive)
	go SendPack([]byte("Hei, Stian"), host, chSend)

	for {
		select {
			case m := <- chReceive:
				fmt.Printf("Received: %s\n", m)
				go SendPack([]byte("Hei, Stian"), host, chSend)
			case <- chSend:
				fmt.Printf("sent!")
		}
	}
}