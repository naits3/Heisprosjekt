package tools

import ("testing")

func TestBroadcastData(t *testing.T) {
	BroadcastElevatorData()
}


func TestBroadcastOrder(t *testing.T) {
	BroadcastOrder()
}

func TestListenData(t *testing.T) {
	ListenForData()
}