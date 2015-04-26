package network

import (
	"encoding/json"
	"Heisprosjekt/src"
)


func PackElevatorData(unpackedElevatorData src.ElevatorData) []byte {
	packedElevatorData,_ := json.Marshal(unpackedElevatorData)
	return packedElevatorData
}

func UnpackElevatorData(packedElevatorData []byte) src.ElevatorData {
	var unpackedElevatorData src.ElevatorData
	json.Unmarshal(packedElevatorData, &unpackedElevatorData)
	return unpackedElevatorData
}

func PackOrder(unpackedOrder src.ButtonOrder) []byte {
	packedOrder, _ := json.Marshal(unpackedOrder)
	return packedOrder
}

func UnpackOrder(packedOrder []byte) src.ButtonOrder {
	var unpackedOrder src.ButtonOrder
	json.Unmarshal(packedOrder, &unpackedOrder)
	return unpackedOrder
}

func PackMessage(unpackedMessage Message) []byte {
	packedMessage, _ := json.Marshal(unpackedMessage)
	return packedMessage
}

func UnpackMessage(packedMessage []byte) Message {
	var unpackedMessage Message
	json.Unmarshal(packedMessage, &unpackedMessage)
	return unpackedMessage
}