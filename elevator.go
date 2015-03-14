package main

import(
	"Heisprosjekt/network"
	//"Heisprosjekt/queue"
)

func main(){
	go network.NetworkHandler()
	for {
		select {
			case data := <- network.ChDataToQueue:
				print(data)
			}
	}
}