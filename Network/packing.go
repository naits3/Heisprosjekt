package network

import ("encoding/json")


func Pack(unpackedMessage elevatorData) []byte {
	packedMessage,_ := json.Marshal(unpackedMessage)
	return packedMessage
}

func Unpack(packedMessage []byte) elevatorData {
	var unpackedMessage elevatorData
	json.Unmarshal(packedMessage, unpackedMessage)
	return unpackedMessage
}