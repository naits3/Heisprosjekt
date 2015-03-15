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
	println(" ID:",queueData.ID)
	println()
}


func PrintQueueArray(queueData []src.ElevatorData) {
	for element := 0; element < len(queueData); element ++ {
		PrintQueue(queueData[element])
	}
}

// TESTED:
func BroadcastElevatorData() {
	host := "192.168.0.255:80"
	addr,_ := net.ResolveUDPAddr("udp",host)
	conn, _ := net.DialUDP("udp",nil,addr)
	var data src.ElevatorData

	for {
		time.Sleep(time.Second)
		conn.Write(network.Pack(data))
	}
}