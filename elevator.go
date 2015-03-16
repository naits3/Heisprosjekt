package main

import(
	"Heisprosjekt/network"
	//"Heisprosjekt/src"
	//"Heisprosjekt/tools"
	"Heisprosjekt/queue"
	"runtime"
)



func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	//var testData src.ElevatorData // FOR TESTING ONLY
	go network.NetworkHandler()
	go queue.QueueHandler()

	// FOR TESTING: ----------------
	// for {
	// 	select {
	// 		case <- network.ChDataToQueue:
	// 			print("Data!")
	// 		case <- network.ChReadyToMerge:
	// 			print("Ready to merge!")
	// 			network.ChQueueReadyToBeSent <- testData
	// 	}
	// }

	for {

	}
}


