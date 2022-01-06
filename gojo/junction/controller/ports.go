package controller

import "../../types"

func getNewPortId(patterns *JoinPatterns, msg types.Packet) {
	messageChannel := types.MessageChannel{
		Ch: make(chan types.Payload, 10),
	}
	(*patterns).firedPorts[(*patterns).portIds] = messageChannel
	msg.Payload.Ch <- (*patterns).portIds
	(*patterns).portIds++
}
