package network

import (
	"encoding/json"
	"Heisprosjekt/src"
)


func Pack(unpackedMessage src.ElevatorData) []byte {
	packedMessage,_ := json.Marshal(unpackedMessage)
	return packedMessage
}

func Unpack(packedMessage []byte) src.ElevatorData {
	var unpackedMessage src.ElevatorData
	json.Unmarshal(packedMessage, &unpackedMessage)
	return unpackedMessage
}