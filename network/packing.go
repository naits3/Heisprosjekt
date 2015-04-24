package network

import (
	"encoding/json"
	"Heisprosjekt/src"
)


func PackQueue(unpackedMessage src.ElevatorData) []byte {
	packedMessage,_ := json.Marshal(unpackedMessage)
	return packedMessage
}

func UnpackQueue(packedMessage []byte) src.ElevatorData {
	var unpackedMessage src.ElevatorData
	json.Unmarshal(packedMessage, &unpackedMessage)
	return unpackedMessage
}

func PackOrder(unpackedMessage src.ButtonOrder) []byte {
	packedMessage, _ := json.Marshal(unpackedMessage)
	return packedMessage
}

func UnpackOrder(packedMessage []byte) src.ButtonOrder {
	var unpackedMessage src.ButtonOrder
	json.Unmarshal(packedMessage, &unpackedMessage)
	return unpackedMessage
}

func PackMessage(unpackedMessage message) []byte {
	packedMessage, _ := json.Marshal(unpackedMessage)
	return packedMessage
}

func UnpackMessage(packedMessage []byte) message {
	var unpackedMessage message
	json.Unmarshal(packedMessage, &unpackedMessage)
	return unpackedMessage
}