package network

import (
	"encoding/json"
)

func PackMessage(unpackedMessage Message) []byte {
	packedMessage, _ := json.Marshal(unpackedMessage)
	return packedMessage
}

func UnpackMessage(packedMessage []byte) Message {
	var unpackedMessage Message
	json.Unmarshal(packedMessage, &unpackedMessage)
	return unpackedMessage
}