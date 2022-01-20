package controller

import "../../types"

func getNewPortId(patterns *JoinPatterns, msg types.Packet) {
	(*patterns).ports[(*patterns).portIds] = make(chan types.Payload)
	msg.Payload.Ch <- (*patterns).portIds
	(*patterns).portIds++
}
