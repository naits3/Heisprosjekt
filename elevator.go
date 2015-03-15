package main

import(
	"Heisprosjekt/network"
	//"Heisprosjekt/src"
	"Heisprosjekt/tools"
	//"Heisprosjekt/queue"
)



func main(){

	go tools.BroadcastElevatorData()
	go network.NetworkHandler()
	
	for {
		select {
			case data := <- network.ChDataToQueue:
				tools.PrintQueue(data)
			case <- network.ChReadyToMerge:
				println("Ready to merge!")			

		}
	}
}


