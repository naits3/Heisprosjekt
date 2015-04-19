package tools

import(
	"Heisprosjekt/src"
	"Heisprosjekt/network"
	"time"
	"net"
)


// TESTED:
func PrintQueue(queueData src.ElevatorData) {
	for row := 0; row < src.N_FLOORS; row ++ {
		for col := 0; col < 2; col ++ {
			print(" ",queueData.OutsideOrders[src.N_FLOORS-1-row][col])
		}
		print(" |",queueData.InsideOrders[src.N_FLOORS-1-row])
		if  row == src.N_FLOORS-1-queueData.Floor {
			print(" *")
			switch queueData.Direction {
				case src.DIR_UP:
					print("U")
				case src.DIR_DOWN:
					print("D")
				default:
					print("I")
			}
			
		}
		println(" ")
	}
	println()
}


func PrintQueueArray(queueData []src.ElevatorData) {
	for element := 0; element < len(queueData); element ++ {
		PrintQueue(queueData[element])
	}
}

// TESTED:
func BroadcastElevatorData() {
	host := "192.168.0.255:20019"
	addr,_ := net.ResolveUDPAddr("udp",host)
	conn, _ := net.DialUDP("udp",nil,addr)
	var data src.ElevatorData

	//data.OutsideOrders[1][1] = 1
	//data.OutsideOrders[2][1] = 1
	//data.OutsideOrders[0][0] = 1
	//data.ID = 143

	for {
		time.Sleep(time.Second)
		conn.Write(network.PackQueue(data))
		println("Wrote data with n = !", len(network.PackQueue(data)))
	}
}

func BroadcastOrder() {
	host := "192.168.0.255:20019"
	addr,_ := net.ResolveUDPAddr("udp",host)
	conn, _ := net.DialUDP("udp",nil,addr)
	data := src.ButtonOrder{3, src.BUTTON_UP}

	for {
		time.Sleep(time.Second)
		conn.Write(network.PackOrder(data))
		println("Wrote data with n = !", len(network.PackOrder(data)))
	}
}


func ListenForData() {
	UDPAddr, _ := net.ResolveUDPAddr("udp",":"+network.PORT)
	var buffer []byte = make([]byte, 1024)
	conn, err := net.ListenUDP("udp", UDPAddr)

	if err != nil {
		println(err)
		return
	}

	defer conn.Close()
	for {
		lengthOfMessage, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			print(err)
			return
		}
				
		
		if (lengthOfMessage < 50) {
			order := network.UnpackOrder(buffer[:lengthOfMessage])
			println("Received order! n = ", lengthOfMessage)
			println("Floor: ", order.Floor, "\nType: ", order.ButtonType)

		}else {
			PrintQueue(network.UnpackQueue(buffer[:lengthOfMessage]))
		}
	}
}


func PrintQueueHandler(mergedQ src.ElevatorData, assignedQ src.ElevatorData) {
	println("The merged Queue: ")
	PrintQueue(mergedQ)

	println("Assigned Queue: ")
	PrintQueue(assignedQ)
}